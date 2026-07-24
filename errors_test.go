package collection

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPErrorClassifications(t *testing.T) {
	tests := []struct {
		status int
		narrow []error
	}{
		{status: 199},
		{status: http.StatusNoContent},
		{status: http.StatusFound},
		{status: http.StatusUnauthorized, narrow: []error{ErrUnauthorized}},
		{status: http.StatusForbidden, narrow: []error{ErrForbidden}},
		{status: http.StatusNotFound, narrow: []error{ErrNotFound, ErrInvalidUserID}},
		{status: http.StatusTooManyRequests, narrow: []error{ErrRateLimited}},
		{status: http.StatusTeapot},
		{status: http.StatusInternalServerError, narrow: []error{ErrServerError}},
	}
	allNarrow := []error{
		ErrUnauthorized,
		ErrForbidden,
		ErrNotFound,
		ErrInvalidUserID,
		ErrRateLimited,
		ErrServerError,
	}
	for _, test := range tests {
		t.Run(http.StatusText(test.status), func(t *testing.T) {
			body := &observedBody{reader: bytes.NewBufferString("upstream-body-secret")}
			client := clientWithResponseBody(body, test.status, http.Header{
				"Location": []string{"https://example.test/location-secret"},
			})
			_, err := client.FetchPage(
				context.Background(),
				"uid-secret",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) || !errors.Is(err, ErrHTTPStatus) {
				t.Fatalf("error = %T %v", err, err)
			}
			if httpErr.StatusCode != test.status || httpErr.Body != "" {
				t.Fatalf("HTTP error = %#v", httpErr)
			}
			for _, candidate := range allNarrow {
				want := containsError(test.narrow, candidate)
				if errors.Is(err, candidate) != want {
					t.Errorf("errors.Is(%v) = %v, want %v", candidate, !want, want)
				}
			}
			for _, marker := range []string{
				"upstream-body-secret",
				"location-secret",
				"uid-secret",
			} {
				if strings.Contains(err.Error(), marker) {
					t.Errorf("error leaked %q: %v", marker, err)
				}
			}
			if len(err.Error()) > 256 || !body.closed.Load() {
				t.Errorf("error length=%d body closed=%v", len(err.Error()), body.closed.Load())
			}
		})
	}
}

func containsError(errorsList []error, target error) bool {
	for _, candidate := range errorsList {
		if candidate == target {
			return true
		}
	}
	return false
}

func TestNilContextIsDirectAndTransportIsNotCalled(t *testing.T) {
	var calls atomic.Int64
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return nil, errors.New("unused")
		})}),
	)
	_, err := client.FetchPage(nil, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	var networkErr *NetworkError
	if !errors.Is(err, ErrNilContext) || errors.As(err, &networkErr) {
		t.Fatalf("error = %T %v", err, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("transport calls = %d", calls.Load())
	}
}

func TestTransportFailureIsSanitized(t *testing.T) {
	const marker = "raw-transport-url-and-secret"
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New(marker)
		})}),
		WithRateLimit(1_000_000, 10),
		WithMaxRetries(0),
	)
	_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	var networkErr *NetworkError
	var retryErr *RetryError
	if !errors.As(err, &networkErr) ||
		!errors.As(err, &retryErr) ||
		!errors.Is(err, ErrTransport) ||
		!errors.Is(err, ErrRetryExhausted) {
		t.Fatalf("error = %T %v", err, err)
	}
	if networkErr.Err != ErrTransport ||
		strings.Contains(err.Error(), marker) ||
		strings.Contains(networkErr.Err.Error(), marker) {
		t.Fatalf("transport detail leaked: retry=%v network=%#v", err, networkErr)
	}
}

func TestParentCancellationAndDeadlineClassification(t *testing.T) {
	tests := []struct {
		name       string
		ctx        func() context.Context
		wantStable error
		wantCtx    error
		timeout    bool
	}{
		{
			name: "canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantStable: ErrCanceled,
			wantCtx:    context.Canceled,
		},
		{
			name: "deadline",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				defer cancel()
				return ctx
			},
			wantStable: ErrTimeout,
			wantCtx:    context.DeadlineExceeded,
			timeout:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(
				"collection-tests/0.1",
				WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					t.Fatal("transport called")
					return nil, nil
				})}),
			)
			_, err := client.FetchPage(
				test.ctx(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			var networkErr *NetworkError
			var retryErr *RetryError
			if !errors.As(err, &networkErr) ||
				networkErr.Timeout != test.timeout ||
				!errors.Is(err, test.wantStable) ||
				!errors.Is(err, test.wantCtx) ||
				errors.As(err, &retryErr) {
				t.Fatalf("error = %T %v, network=%#v", err, err, networkErr)
			}
		})
	}
}

func TestDecodeProtocolAndRetryErrorsSupportIsAndAs(t *testing.T) {
	decodeErr := newDecodeError("test-secret")
	var decoded *DecodeError
	if !errors.As(decodeErr, &decoded) || !errors.Is(decodeErr, ErrDecode) {
		t.Fatalf("decode error = %v", decodeErr)
	}
	protocolErr := newProtocolError()
	var protocol *ProtocolError
	if !errors.As(protocolErr, &protocol) || !errors.Is(protocolErr, ErrProtocol) {
		t.Fatalf("protocol error = %v", protocolErr)
	}
	retryErr := &RetryError{Attempts: 2, Err: &HTTPError{StatusCode: http.StatusTooManyRequests}}
	var retry *RetryError
	var httpErr *HTTPError
	if !errors.As(retryErr, &retry) ||
		!errors.As(retryErr, &httpErr) ||
		!errors.Is(retryErr, ErrRetryExhausted) ||
		!errors.Is(retryErr, ErrRateLimited) {
		t.Fatalf("retry error = %v", retryErr)
	}
}

type contextBody struct {
	ctx    context.Context
	closed atomic.Bool
}

func (body *contextBody) Read([]byte) (int, error) {
	<-body.ctx.Done()
	return 0, body.ctx.Err()
}

func (body *contextBody) Close() error {
	body.closed.Store(true)
	return nil
}

func TestParentCancellationDuringBodyReadWinsAndCloses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var body *contextBody
	started := make(chan struct{})
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			body = &contextBody{ctx: request.Context()}
			close(started)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       body,
			}, nil
		})}),
		WithRateLimit(1_000_000, 10),
		WithMaxRetries(3),
	)
	result := make(chan error, 1)
	go func() {
		_, err := client.FetchPage(ctx, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
		result <- err
	}()
	<-started
	cancel()
	err := <-result
	var networkErr *NetworkError
	if !errors.As(err, &networkErr) ||
		!errors.Is(err, ErrCanceled) ||
		!errors.Is(err, context.Canceled) {
		t.Fatalf("error = %T %v", err, err)
	}
	if body == nil || !body.closed.Load() {
		t.Fatal("body was not closed")
	}
}

var _ io.ReadCloser = (*contextBody)(nil)
