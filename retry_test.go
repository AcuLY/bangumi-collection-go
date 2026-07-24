package collection

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func successResponse(limit, offset int) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(fixturePage(0, limit, offset))),
	}
}

func statusResponse(status int, retryAfter string) *http.Response {
	header := make(http.Header)
	if retryAfter != "" {
		header.Set("Retry-After", retryAfter)
	}
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(bytes.NewBufferString("discarded-secret")),
	}
}

func TestRetryableMatrixStopsOnSuccess(t *testing.T) {
	tests := []struct {
		name      string
		transport func(*atomic.Int64) roundTripFunc
	}{
		{
			name: "transport",
			transport: func(calls *atomic.Int64) roundTripFunc {
				return func(*http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						return nil, errors.New("raw transport secret")
					}
					return successResponse(1, 0), nil
				}
			},
		},
		{
			name: "attempt timeout",
			transport: func(calls *atomic.Int64) roundTripFunc {
				return func(request *http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						<-request.Context().Done()
						return nil, request.Context().Err()
					}
					return successResponse(1, 0), nil
				}
			},
		},
		{
			name: "rate limited",
			transport: func(calls *atomic.Int64) roundTripFunc {
				return func(*http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						return statusResponse(http.StatusTooManyRequests, ""), nil
					}
					return successResponse(1, 0), nil
				}
			},
		},
		{
			name: "server error",
			transport: func(calls *atomic.Int64) roundTripFunc {
				return func(*http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						return statusResponse(http.StatusBadGateway, ""), nil
					}
					return successResponse(1, 0), nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int64
			client := NewClient(
				"collection-tests/0.1",
				WithHTTPClient(&http.Client{Transport: test.transport(&calls)}),
				WithRequestTimeout(5*time.Millisecond),
				WithMaxRetries(2),
				WithRetryInterval(time.Nanosecond),
				WithRateLimit(1_000_000, 10),
			)
			client.sleep = func(context.Context, time.Duration) error { return nil }
			page, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			if err != nil {
				t.Fatal(err)
			}
			if page == nil || calls.Load() != 2 {
				t.Fatalf("page = %#v, calls = %d", page, calls.Load())
			}
		})
	}
}

func TestTerminalResultsAreNotRetried(t *testing.T) {
	tests := []struct {
		name     string
		response func() *http.Response
	}{
		{name: "informational", response: func() *http.Response { return statusResponse(199, "") }},
		{name: "no content", response: func() *http.Response { return statusResponse(204, "") }},
		{name: "redirect", response: func() *http.Response { return statusResponse(302, "") }},
		{name: "unauthorized", response: func() *http.Response { return statusResponse(401, "") }},
		{name: "forbidden", response: func() *http.Response { return statusResponse(403, "") }},
		{name: "not found", response: func() *http.Response { return statusResponse(404, "") }},
		{name: "other client", response: func() *http.Response { return statusResponse(418, "") }},
		{name: "decode", response: func() *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewBufferString("not-json")),
			}
		}},
		{name: "protocol", response: func() *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(fixturePage(0, 2, 0))),
			}
		}},
		{name: "oversized", response: func() *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(bytes.Repeat([]byte("x"), maxSuccessBodyBytes+1))),
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int64
			client := NewClient(
				"collection-tests/0.1",
				WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					calls.Add(1)
					return test.response(), nil
				})}),
				WithMaxRetries(3),
				WithRateLimit(1_000_000, 10),
			)
			_, _ = client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			if calls.Load() != 1 {
				t.Fatalf("transport calls = %d", calls.Load())
			}
		})
	}
}

func TestRetryAfterAndJitterUseInjectedHooks(t *testing.T) {
	fixedNow := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		header    string
		random    float64
		maxDelay  time.Duration
		wantDelay time.Duration
	}{
		{name: "delta", header: "2", random: 0, maxDelay: 5 * time.Second, wantDelay: 2 * time.Second},
		{
			name:      "date",
			header:    fixedNow.Add(3 * time.Second).Format(http.TimeFormat),
			random:    0,
			maxDelay:  5 * time.Second,
			wantDelay: 3 * time.Second,
		},
		{name: "malformed uses jitter", header: "later", random: 0.5, maxDelay: 5 * time.Second, wantDelay: 500 * time.Millisecond},
		{name: "signed delta is malformed", header: "+2", random: 0.5, maxDelay: 5 * time.Second, wantDelay: 500 * time.Millisecond},
		{
			name:      "past uses jitter",
			header:    fixedNow.Add(-time.Second).Format(http.TimeFormat),
			random:    0.5,
			maxDelay:  5 * time.Second,
			wantDelay: 500 * time.Millisecond,
		},
		{name: "negative uses jitter", header: "-3", random: 0.5, maxDelay: 5 * time.Second, wantDelay: 500 * time.Millisecond},
		{name: "server delay capped", header: "20", random: 0, maxDelay: 1500 * time.Millisecond, wantDelay: 1500 * time.Millisecond},
		{
			name:      "arbitrarily long delta saturates",
			header:    strings.Repeat("9", 1000),
			random:    0,
			maxDelay:  5 * time.Second,
			wantDelay: 5 * time.Second,
		},
		{
			name:      "arbitrarily long leading zeros preserve value",
			header:    strings.Repeat("0", 1000) + "2",
			random:    0,
			maxDelay:  5 * time.Second,
			wantDelay: 2 * time.Second,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int64
			var delaysMu sync.Mutex
			var delays []time.Duration
			client := NewClient(
				"collection-tests/0.1",
				WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						return statusResponse(http.StatusServiceUnavailable, test.header), nil
					}
					return successResponse(1, 0), nil
				})}),
				WithMaxRetries(1),
				WithRetryInterval(time.Second),
				WithMaxRetryDelay(test.maxDelay),
				WithRateLimit(1_000_000, 10),
			)
			client.now = func() time.Time { return fixedNow }
			client.randomFloat = func() float64 { return test.random }
			client.sleep = func(_ context.Context, delay time.Duration) error {
				delaysMu.Lock()
				delays = append(delays, delay)
				delaysMu.Unlock()
				return nil
			}
			if _, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			); err != nil {
				t.Fatal(err)
			}
			delaysMu.Lock()
			defer delaysMu.Unlock()
			if len(delays) != 1 || delays[0] != test.wantDelay {
				t.Fatalf("delays = %v, want %s", delays, test.wantDelay)
			}
		})
	}
}

func TestFullJitterBoundsAndOverflowSafeCap(t *testing.T) {
	client := NewClient(
		"collection-tests/0.1",
		WithRetryInterval(time.Second),
		WithMaxRetryDelay(3*time.Second),
	)
	tests := []struct {
		random float64
		retry  int
		want   time.Duration
	}{
		{random: -1, retry: 1, want: 0},
		{random: 0, retry: 1, want: 0},
		{random: 0.5, retry: 1, want: 500 * time.Millisecond},
		{random: 1, retry: 1, want: time.Second},
		{random: 1, retry: 2, want: 2 * time.Second},
		{random: 1, retry: 3, want: 3 * time.Second},
		{random: 1, retry: 1000, want: 3 * time.Second},
		{random: 2, retry: 1000, want: 3 * time.Second},
	}
	for _, test := range tests {
		client.randomFloat = func() float64 { return test.random }
		if delay := client.retryDelay(test.retry, 0); delay != test.want {
			t.Errorf("random=%v retry=%d: delay=%s, want %s", test.random, test.retry, delay, test.want)
		}
	}

	maxDurationClient := NewClient(
		"collection-tests/0.1",
		WithRetryInterval(time.Duration(1<<63-1)),
		WithMaxRetryDelay(time.Duration(1<<63-1)),
	)
	maxDurationClient.randomFloat = func() float64 { return 1 }
	if delay := maxDurationClient.retryDelay(1000, 0); delay != time.Duration(1<<63-1) {
		t.Fatalf("maximum-duration jitter overflowed: %s", delay)
	}
}

func TestRetryExhaustionPreservesAttemptsAndLastClassification(t *testing.T) {
	var calls atomic.Int64
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return statusResponse(http.StatusTooManyRequests, ""), nil
		})}),
		WithMaxRetries(2),
		WithRetryInterval(time.Nanosecond),
		WithRateLimit(1_000_000, 10),
	)
	client.sleep = func(context.Context, time.Duration) error { return nil }

	_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	var retryErr *RetryError
	var httpErr *HTTPError
	if !errors.As(err, &retryErr) ||
		!errors.As(err, &httpErr) ||
		retryErr.Attempts != 3 ||
		!errors.Is(err, ErrRetryExhausted) ||
		!errors.Is(err, ErrRateLimited) ||
		calls.Load() != 3 {
		t.Fatalf("error=%T %v retry=%#v calls=%d", err, err, retryErr, calls.Load())
	}
}

func TestCancellationDuringRetryWaitIsTerminal(t *testing.T) {
	var calls atomic.Int64
	sleepStarted := make(chan struct{})
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return statusResponse(http.StatusServiceUnavailable, ""), nil
		})}),
		WithMaxRetries(3),
		WithRateLimit(1_000_000, 10),
	)
	client.sleep = func(ctx context.Context, _ time.Duration) error {
		close(sleepStarted)
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := client.FetchPage(ctx, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
		result <- err
	}()
	<-sleepStarted
	cancel()
	err := <-result
	var retryErr *RetryError
	if !errors.Is(err, ErrCanceled) ||
		!errors.Is(err, context.Canceled) ||
		errors.As(err, &retryErr) ||
		calls.Load() != 1 {
		t.Fatalf("error=%T %v calls=%d", err, err, calls.Load())
	}
}

func TestEveryRetryParticipatesInSharedLimiters(t *testing.T) {
	var calls atomic.Int64
	var badPermit atomic.Bool
	var client *Client
	client = NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if len(client.inflight) != 1 {
				badPermit.Store(true)
			}
			if calls.Add(1) == 1 {
				return statusResponse(http.StatusBadGateway, ""), nil
			}
			limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
			offset, _ := strconv.Atoi(request.URL.Query().Get("offset"))
			return successResponse(limit, offset), nil
		})}),
		WithMaxRetries(1),
		WithRateLimit(0.000001, 2),
		WithConcurrencyLimit(1),
	)
	client.sleep = func(context.Context, time.Duration) error { return nil }

	_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 || badPermit.Load() || len(client.inflight) != 0 {
		t.Fatalf("calls=%d bad permit=%v remaining permits=%d", calls.Load(), badPermit.Load(), len(client.inflight))
	}
	if tokens := client.requestLimiter.Tokens(); tokens > 0.01 {
		t.Fatalf("retry did not consume rate token: %f", tokens)
	}
}
