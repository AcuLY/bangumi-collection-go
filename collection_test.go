package collection

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func requestPageCoordinates(t *testing.T, request *http.Request) (CollectionType, int, int) {
	t.Helper()
	collectionValue, err := strconv.Atoi(request.URL.Query().Get("type"))
	if err != nil {
		t.Errorf("collection type: %v", err)
	}
	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		t.Errorf("limit: %v", err)
	}
	offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
	if err != nil {
		t.Errorf("offset: %v", err)
	}
	return CollectionType(collectionValue), limit, offset
}

func TestFetchEmptySuccessAndFetchPageOrderClamp(t *testing.T) {
	t.Run("empty aggregate", func(t *testing.T) {
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, limit, offset := requestPageCoordinates(t, request)
			if limit != pageSize || offset != 0 {
				t.Errorf("limit=%d offset=%d", limit, offset)
			}
			writeJSON(t, writer, fixturePage(0, limit, offset))
		}))
		subjects, err := client.Fetch(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
		)
		if err != nil {
			t.Fatal(err)
		}
		if subjects == nil || len(subjects) != 0 {
			t.Fatalf("subjects = %#v", subjects)
		}
	})

	t.Run("page clamp and upstream order", func(t *testing.T) {
		var calls atomic.Int64
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			calls.Add(1)
			collectionType, limit, offset := requestPageCoordinates(t, request)
			if limit != 50 || offset != 0 {
				t.Errorf("limit=%d offset=%d", limit, offset)
			}
			writeJSON(t, writer, fixturePage(
				2,
				limit,
				offset,
				fixtureCollection(2, SubjectTypeAnime, collectionType),
				fixtureCollection(1, SubjectTypeAnime, collectionType),
			))
		}))
		page, err := client.FetchPage(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
			100,
			-5,
		)
		if err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 1 ||
			len(page.Data) != 2 ||
			page.Data[0].SubjectID != 2 ||
			page.Data[1].SubjectID != 1 {
			t.Fatalf("page = %#v, calls=%d", page, calls.Load())
		}
	})
}

func TestFetchInputValidationHappensBeforeTransport(t *testing.T) {
	var calls atomic.Int64
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return nil, errors.New("must not run")
		})}),
	)
	tests := []struct {
		name string
		run  func() error
		want error
	}{
		{
			name: "invalid subject",
			run: func() error {
				_, err := client.Fetch(context.Background(), "uid", SubjectType(5), CollectionTypeDone)
				return err
			},
			want: ErrInvalidSubjectType,
		},
		{
			name: "invalid collection",
			run: func() error {
				_, err := client.Fetch(context.Background(), "uid", SubjectTypeAnime, CollectionType(6))
				return err
			},
			want: ErrInvalidCollectionType,
		},
		{
			name: "no collection",
			run: func() error {
				_, err := client.Fetch(context.Background(), "uid", SubjectTypeAnime)
				return err
			},
			want: ErrNoCollectionTypes,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("transport calls = %d", calls.Load())
	}
}

func TestFetchMultiplePagesTypesAndCanonicalOrder(t *testing.T) {
	type requestKey struct {
		collectionType CollectionType
		offset         int
	}
	var mu sync.Mutex
	counts := make(map[requestKey]int)
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			if limit != pageSize {
				t.Errorf("observed old probe limit %d", limit)
			}
			key := requestKey{collectionType: collectionType, offset: offset}
			mu.Lock()
			counts[key]++
			mu.Unlock()

			if offset == 50 {
				if collectionType == CollectionTypeWish {
					time.Sleep(15 * time.Millisecond)
				} else {
					time.Sleep(time.Millisecond)
				}
			}
			var item map[string]any
			switch key {
			case requestKey{collectionType: CollectionTypeWish, offset: 0}:
				item = fixtureCollection(3, SubjectTypeAnime, collectionType)
			case requestKey{collectionType: CollectionTypeWish, offset: 50}:
				item = fixtureCollection(1, SubjectTypeAnime, collectionType)
			case requestKey{collectionType: CollectionTypeDone, offset: 0}:
				item = fixtureCollection(1, SubjectTypeAnime, collectionType)
			case requestKey{collectionType: CollectionTypeDone, offset: 50}:
				item = fixtureCollection(2, SubjectTypeAnime, collectionType)
			default:
				t.Errorf("unexpected request %#v", key)
				writeJSON(t, writer, fixturePage(0, limit, offset))
				return
			}
			writeJSON(t, writer, fixturePage(51, limit, offset, item))
		}),
		WithConcurrencyLimit(2),
	)

	subjects, err := client.Fetch(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
		CollectionTypeWish,
		CollectionTypeDone,
	)
	if err != nil {
		t.Fatal(err)
	}
	gotKeys := make([]subjectKey, 0, len(subjects))
	for _, subject := range subjects {
		gotKeys = append(gotKeys, subjectKey{
			subjectType: subject.SubjectType,
			subjectID:   subject.SubjectID,
			collection:  subject.Type,
		})
	}
	wantKeys := []subjectKey{
		{subjectType: SubjectTypeAnime, subjectID: 1, collection: CollectionTypeWish},
		{subjectType: SubjectTypeAnime, subjectID: 1, collection: CollectionTypeDone},
		{subjectType: SubjectTypeAnime, subjectID: 2, collection: CollectionTypeDone},
		{subjectType: SubjectTypeAnime, subjectID: 3, collection: CollectionTypeWish},
	}
	if !reflect.DeepEqual(gotKeys, wantKeys) {
		t.Fatalf("keys = %#v, want %#v", gotKeys, wantKeys)
	}

	mu.Lock()
	defer mu.Unlock()
	wantRequests := []requestKey{
		{collectionType: CollectionTypeWish, offset: 0},
		{collectionType: CollectionTypeWish, offset: 50},
		{collectionType: CollectionTypeDone, offset: 0},
		{collectionType: CollectionTypeDone, offset: 50},
	}
	for _, key := range wantRequests {
		if counts[key] != 1 {
			t.Errorf("request %#v count = %d", key, counts[key])
		}
	}
	if len(counts) != len(wantRequests) {
		t.Fatalf("request counts = %#v", counts)
	}
}

func TestFetchFirstPageFixesMovingTotal(t *testing.T) {
	var mu sync.Mutex
	var offsets []int
	client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, limit, offset := requestPageCoordinates(t, request)
		mu.Lock()
		offsets = append(offsets, offset)
		mu.Unlock()
		total := 101
		switch offset {
		case 50:
			total = 1_000_000
		case 100:
			total = 0
		}
		writeJSON(t, writer, fixturePage(total, limit, offset))
	}))

	subjects, err := client.Fetch(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
	)
	if err != nil {
		t.Fatal(err)
	}
	if subjects == nil || len(subjects) != 0 {
		t.Fatalf("subjects = %#v", subjects)
	}
	mu.Lock()
	defer mu.Unlock()
	sort.Ints(offsets)
	if !reflect.DeepEqual(offsets, []int{0, 50, 100}) {
		t.Fatalf("offsets = %v", offsets)
	}
}

func TestFetchDuplicateWinnerAndCollectionTypePreservation(t *testing.T) {
	client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		collectionType, limit, offset := requestPageCoordinates(t, request)
		switch {
		case collectionType == CollectionTypeWish:
			item := fixtureCollection(2, SubjectTypeAnime, collectionType)
			item["comment"] = "wish"
			writeJSON(t, writer, fixturePage(1, limit, offset, item))
		case offset == 0:
			first := fixtureCollection(2, SubjectTypeAnime, collectionType)
			first["comment"] = "smallest-coordinate"
			second := fixtureCollection(2, SubjectTypeAnime, collectionType)
			second["comment"] = "same-page-later"
			writeJSON(t, writer, fixturePage(51, limit, offset, first, second))
		case offset == 50:
			later := fixtureCollection(2, SubjectTypeAnime, collectionType)
			later["comment"] = "later-page"
			other := fixtureCollection(1, SubjectTypeAnime, collectionType)
			writeJSON(t, writer, fixturePage(51, limit, offset, later, other))
		default:
			t.Errorf("unexpected offset %d", offset)
		}
	}))

	subjects, err := client.Fetch(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
		CollectionTypeWish,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(subjects) != 3 {
		t.Fatalf("subjects = %#v", subjects)
	}
	if subjects[0].SubjectID != 1 ||
		subjects[1].SubjectID != 2 ||
		subjects[1].Type != CollectionTypeWish ||
		subjects[2].SubjectID != 2 ||
		subjects[2].Type != CollectionTypeDone ||
		subjects[2].Comment != "smallest-coordinate" {
		t.Fatalf("subjects = %#v", subjects)
	}
}

func TestFetchProtocolBoundariesStopBeforeScheduling(t *testing.T) {
	tests := []struct {
		name string
		body func() []byte
	}{
		{name: "total above maximum", body: func() []byte { return fixturePage(1_000_001, 50, 0) }},
		{name: "near max int", body: func() []byte { return fixturePage(maxInt()-1, 50, 0) }},
		{name: "mismatched limit", body: func() []byte { return fixturePage(0, 49, 0) }},
		{name: "mismatched offset", body: func() []byte { return fixturePage(0, 50, 1) }},
		{name: "first page data exceeds zero total", body: func() []byte {
			return fixturePage(
				0,
				50,
				0,
				fixtureCollection(1, SubjectTypeAnime, CollectionTypeDone),
			)
		}},
		{name: "first page data exceeds positive total", body: func() []byte {
			return fixturePage(
				1,
				50,
				0,
				fixtureCollection(1, SubjectTypeAnime, CollectionTypeDone),
				fixtureCollection(2, SubjectTypeAnime, CollectionTypeDone),
			)
		}},
		{name: "too much data", body: func() []byte {
			items := make([]map[string]any, 51)
			for index := range items {
				items[index] = fixtureCollection(index+1, SubjectTypeAnime, CollectionTypeDone)
			}
			return fixturePage(51, 50, 0, items...)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int64
			client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				writeJSON(t, writer, test.body())
			}))
			subjects, err := client.Fetch(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
			)
			if !errors.Is(err, ErrProtocol) || subjects != nil {
				t.Fatalf("subjects=%#v error=%T %v", subjects, err, err)
			}
			if calls.Load() != 1 {
				t.Fatalf("transport calls = %d", calls.Load())
			}
		})
	}
}

func TestFetchMaximumTotalUsesBoundedWorkers(t *testing.T) {
	var active atomic.Int64
	var peak atomic.Int64
	started := make(chan struct{}, 3)
	finished := make(chan struct{}, 3)
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, limit, offset := requestPageCoordinates(t, request)
			if offset == 0 {
				writeJSON(t, writer, fixturePage(maxPageTotal, limit, offset))
				return
			}
			current := active.Add(1)
			for {
				old := peak.Load()
				if current <= old || peak.CompareAndSwap(old, current) {
					break
				}
			}
			started <- struct{}{}
			<-request.Context().Done()
			active.Add(-1)
			finished <- struct{}{}
		}),
		WithConcurrencyLimit(3),
		WithRequestTimeout(time.Second),
	)

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := client.Fetch(ctx, "uid", SubjectTypeAnime, CollectionTypeDone)
		result <- err
	}()
	for range 3 {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("bounded workers did not start")
		}
	}
	if peak.Load() > 3 {
		t.Fatalf("peak workers = %d", peak.Load())
	}
	cancel()
	err := <-result
	if !errors.Is(err, ErrCanceled) || !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %T %v", err, err)
	}
	for range 3 {
		select {
		case <-finished:
		case <-time.After(time.Second):
			t.Fatal("canceled server request did not finish")
		}
	}
	if active.Load() != 0 || peak.Load() != 3 || len(client.inflight) != 0 {
		t.Fatalf("active=%d peak=%d permits=%d", active.Load(), peak.Load(), len(client.inflight))
	}
}

func TestFetchRemainingPageWorkerPeak(t *testing.T) {
	var active atomic.Int64
	var peak atomic.Int64
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, limit, offset := requestPageCoordinates(t, request)
			if offset > 0 {
				current := active.Add(1)
				for {
					old := peak.Load()
					if current <= old || peak.CompareAndSwap(old, current) {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
				active.Add(-1)
			}
			writeJSON(t, writer, fixturePage(151, limit, offset))
		}),
		WithConcurrencyLimit(2),
	)
	if _, err := client.Fetch(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
	); err != nil {
		t.Fatal(err)
	}
	if peak.Load() != 2 || active.Load() != 0 {
		t.Fatalf("peak=%d active=%d", peak.Load(), active.Load())
	}
}

func TestFetchTerminalPageCancelsSiblingAndReturnsNoPartial(t *testing.T) {
	var remainingStarted atomic.Int64
	bothReady := make(chan struct{})
	siblingCanceled := make(chan struct{})
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			if offset == 0 {
				writeJSON(t, writer, fixturePage(
					101,
					limit,
					offset,
					fixtureCollection(1, SubjectTypeAnime, collectionType),
				))
				return
			}
			if remainingStarted.Add(1) == 2 {
				close(bothReady)
			}
			<-bothReady
			if offset == 50 {
				writer.WriteHeader(http.StatusTeapot)
				return
			}
			<-request.Context().Done()
			close(siblingCanceled)
		}),
		WithConcurrencyLimit(2),
		WithMaxRetries(3),
	)

	subjects, err := client.Fetch(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
	)
	var httpErr *HTTPError
	if subjects != nil || !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTeapot {
		t.Fatalf("subjects=%#v error=%T %v", subjects, err, err)
	}
	select {
	case <-siblingCanceled:
	case <-time.After(time.Second):
		t.Fatal("sibling page did not observe cancellation")
	}
	if len(client.inflight) != 0 {
		t.Fatalf("permit leak = %d", len(client.inflight))
	}
}

func TestFetchRepeatedRunDeterminism(t *testing.T) {
	var calls atomic.Int64
	client, _ := newLoopbackServerClient(
		t,
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			if offset == 50 && calls.Add(1)%2 == 0 {
				time.Sleep(2 * time.Millisecond)
			}
			id := 3
			if offset == 50 {
				id = 1
			} else if collectionType == CollectionTypeDone {
				id = 2
			}
			writeJSON(t, writer, fixturePage(
				51,
				limit,
				offset,
				fixtureCollection(id, SubjectTypeAnime, collectionType),
			))
		}),
		WithConcurrencyLimit(3),
	)

	var first []*Subject
	for iteration := 0; iteration < 5; iteration++ {
		subjects, err := client.Fetch(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
			CollectionTypeWish,
			CollectionTypeDone,
		)
		if err != nil {
			t.Fatal(err)
		}
		if iteration == 0 {
			first = subjects
			continue
		}
		if !reflect.DeepEqual(subjects, first) {
			t.Fatalf("iteration %d differs:\nfirst=%#v\nlater=%#v", iteration, first, subjects)
		}
	}
}
