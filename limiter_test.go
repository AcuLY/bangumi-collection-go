package collection

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func writeEmptyPageForRequest(t *testing.T, writer http.ResponseWriter, request *http.Request) {
	t.Helper()
	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		t.Errorf("limit: %v", err)
		return
	}
	offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
	if err != nil {
		t.Errorf("offset: %v", err)
		return
	}
	writeJSON(t, writer, fixturePage(0, limit, offset))
}

func TestClientWideRateLimitAcrossConcurrentCalls(t *testing.T) {
	var mu sync.Mutex
	var starts []time.Time
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			mu.Lock()
			starts = append(starts, time.Now())
			mu.Unlock()
			writeEmptyPageForRequest(t, writer, request)
		}),
		WithRateLimit(25, 1),
	)

	start := make(chan struct{})
	errs := make(chan error, 3)
	for range 3 {
		go func() {
			<-start
			_, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			errs <- err
		}()
	}
	close(start)
	for range 3 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}

	mu.Lock()
	defer mu.Unlock()
	if len(starts) != 3 {
		t.Fatalf("starts = %v", starts)
	}
	if span := starts[2].Sub(starts[0]); span < 65*time.Millisecond {
		t.Fatalf("three starts spanned only %s", span)
	}
}

func TestClientWideInFlightLimitAndCrossOperationSharing(t *testing.T) {
	var active atomic.Int64
	var peak atomic.Int64
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			current := active.Add(1)
			for {
				old := peak.Load()
				if current <= old || peak.CompareAndSwap(old, current) {
					break
				}
			}
			started <- struct{}{}
			<-release
			active.Add(-1)
			writeEmptyPageForRequest(t, writer, request)
		}),
		WithConcurrencyLimit(1),
		WithRateLimit(1_000_000, 10),
	)

	errs := make(chan error, 2)
	go func() {
		_, err := client.Fetch(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
		)
		errs <- err
	}()
	go func() {
		_, err := client.FetchPage(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeWish,
			1,
			0,
		)
		errs <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("first request did not start")
	}
	select {
	case <-started:
		t.Fatal("second request exceeded in-flight limit")
	case <-time.After(30 * time.Millisecond):
	}
	release <- struct{}{}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("second request did not start after release")
	}
	release <- struct{}{}

	for range 2 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	if peak.Load() != 1 || active.Load() != 0 || len(client.inflight) != 0 {
		t.Fatalf("peak=%d active=%d permits=%d", peak.Load(), active.Load(), len(client.inflight))
	}
}

func TestSeparateClientsHaveIndependentLimits(t *testing.T) {
	var active atomic.Int64
	var peak atomic.Int64
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		current := active.Add(1)
		for {
			old := peak.Load()
			if current <= old || peak.CompareAndSwap(old, current) {
				break
			}
		}
		started <- struct{}{}
		<-release
		active.Add(-1)
		writeEmptyPageForRequest(t, writer, request)
	})

	first, server := newLoopbackServerClient(
		t,
		handler,
		WithConcurrencyLimit(1),
		WithRateLimit(1, 1),
	)
	second := NewClient(
		"collection-tests/0.1",
		WithEndpoint(server.URL),
		WithHTTPClient(&http.Client{
			Transport: loopbackOnlyTransport{base: server.Client().Transport},
		}),
		WithConcurrencyLimit(1),
		WithRateLimit(1, 1),
		WithMaxRetries(0),
	)

	errs := make(chan error, 2)
	for _, client := range []*Client{first, second} {
		go func(client *Client) {
			_, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			errs <- err
		}(client)
	}
	for range 2 {
		select {
		case <-started:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("clients did not start independently")
		}
	}
	release <- struct{}{}
	release <- struct{}{}
	for range 2 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	if peak.Load() != 2 {
		t.Fatalf("peak active = %d", peak.Load())
	}
}

func TestRateWaitWouldExceedParentDeadlineIsTerminal(t *testing.T) {
	var calls atomic.Int64
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			calls.Add(1)
			writeEmptyPageForRequest(t, writer, request)
		}),
		WithRateLimit(0.01, 1),
		WithMaxRetries(3),
	)
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := client.FetchPage(ctx, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	if !errors.Is(err, ErrTimeout) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %T %v", err, err)
	}
	var networkErr *NetworkError
	var retryErr *RetryError
	if !errors.As(err, &networkErr) ||
		errors.As(err, &retryErr) ||
		errors.Is(err, ErrRetryExhausted) {
		t.Fatalf("rate deadline classification = %T %v", err, err)
	}
	if ctx.Err() != nil {
		t.Fatalf("rate wait returned after parent deadline: %v", ctx.Err())
	}
	if calls.Load() != 1 || len(client.inflight) != 0 {
		t.Fatalf("transport calls = %d, permits = %d", calls.Load(), len(client.inflight))
	}
}

func TestCancellationWhileSemaphoreWaitingBalancesPermit(t *testing.T) {
	var calls atomic.Int64
	started := make(chan struct{})
	release := make(chan struct{})
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			call := calls.Add(1)
			if call == 1 {
				close(started)
				<-release
			}
			writeEmptyPageForRequest(t, writer, request)
		}),
		WithConcurrencyLimit(1),
		WithRateLimit(1_000_000, 10),
	)

	firstDone := make(chan error, 1)
	go func() {
		_, err := client.FetchPage(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
			1,
			0,
		)
		firstDone <- err
	}()
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := client.FetchPage(ctx, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	if !errors.Is(err, ErrTimeout) || calls.Load() != 1 {
		t.Fatalf("error = %v, calls = %d", err, calls.Load())
	}
	close(release)
	if err := <-firstDone; err != nil {
		t.Fatal(err)
	}
	if len(client.inflight) != 0 {
		t.Fatalf("permit leak = %d", len(client.inflight))
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
	if calls.Load() != 2 || len(client.inflight) != 0 {
		t.Fatalf("calls = %d, permits = %d", calls.Load(), len(client.inflight))
	}
}
