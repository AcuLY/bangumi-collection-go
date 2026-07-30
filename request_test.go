package collection

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAnonymousRequestConstructionAndUIDEscaping(t *testing.T) {
	const suppliedUID = "\u2003Mixed/Name?Value\u2003"
	const normalizedUID = "Mixed/Name?Value"
	var observed atomic.Bool

	client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		observed.Store(true)
		if request.Method != http.MethodGet {
			t.Errorf("method = %s", request.Method)
		}
		expectedPath := "/v0/users/" + url.PathEscape(normalizedUID) + "/collections"
		if request.URL.EscapedPath() != expectedPath {
			t.Errorf("escaped path = %q, want %q", request.URL.EscapedPath(), expectedPath)
		}
		query := request.URL.Query()
		if query.Get("subject_type") != "2" ||
			query.Get("type") != "2" ||
			query.Get("limit") != "50" ||
			query.Get("offset") != "0" ||
			len(query) != 4 {
			t.Errorf("query = %#v", query)
		}
		if request.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept = %q", request.Header.Get("Accept"))
		}
		if request.Header.Get("Authorization") != "" || request.Header.Get("Cookie") != "" {
			t.Errorf("credential headers: Authorization=%q Cookie=%q",
				request.Header.Get("Authorization"),
				request.Header.Get("Cookie"))
		}
		writeJSON(t, writer, fixturePage(0, 50, 0))
	}))

	result, err := client.FetchPage(
		context.Background(),
		suppliedUID,
		SubjectTypeAnime,
		CollectionTypeDone,
		50,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !observed.Load() || result.Data == nil || len(result.Data) != 0 {
		t.Fatalf("result = %#v, observed = %v", result, observed.Load())
	}
}

func TestNamedSubjectConstantsConstructOfficialQueryValues(t *testing.T) {
	tests := []struct {
		name        string
		subjectType SubjectType
		want        string
	}{
		{name: "music", subjectType: SubjectTypeMusic, want: "3"},
		{name: "game", subjectType: SubjectTypeGame, want: "4"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var observed string
			client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				observed = request.URL.Query().Get("subject_type")
				writeJSON(t, writer, fixturePage(0, 1, 0))
			}))

			if _, err := client.FetchPage(
				context.Background(),
				"uid",
				test.subjectType,
				CollectionTypeDone,
				1,
				0,
			); err != nil {
				t.Fatal(err)
			}
			if observed != test.want {
				t.Fatalf("subject_type = %q, want %q", observed, test.want)
			}
		})
	}
}

func TestDotOnlyUIDRemainsOneOpaqueEscapedSegment(t *testing.T) {
	tests := []struct {
		uid     string
		escaped string
	}{
		{uid: ".", escaped: "%2E"},
		{uid: "..", escaped: "%2E%2E"},
	}
	for _, test := range tests {
		t.Run(test.uid, func(t *testing.T) {
			var observed atomic.Bool
			mux := http.NewServeMux()
			mux.HandleFunc("GET /v0/users/{uid}/collections", func(writer http.ResponseWriter, request *http.Request) {
				observed.Store(true)
				if request.PathValue("uid") != test.uid {
					t.Errorf("path value = %q, want %q", request.PathValue("uid"), test.uid)
				}
				requestPath := strings.SplitN(request.RequestURI, "?", 2)[0]
				wantPath := "/v0/users/" + test.escaped + "/collections"
				if requestPath != wantPath {
					t.Errorf("RequestURI path = %q, want %q", requestPath, wantPath)
				}
				segments := strings.Split(strings.TrimPrefix(requestPath, "/"), "/")
				if len(segments) != 4 || segments[2] != test.escaped {
					t.Errorf("opaque segments = %#v", segments)
				}
				writeJSON(t, writer, fixturePage(0, 1, 0))
			})

			client, _ := newLoopbackServerClient(t, mux)
			page, err := client.FetchPage(
				context.Background(),
				test.uid,
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			if err != nil {
				t.Fatal(err)
			}
			if !observed.Load() || page.Data == nil {
				t.Fatalf("page = %#v, observed = %v", page, observed.Load())
			}
		})
	}
}

func TestUIDValidationFailsBeforeTransport(t *testing.T) {
	var calls atomic.Int64
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return nil, errors.New("must not run")
		})}),
	)
	invalidUTF8 := string([]byte{0xff})
	tests := []struct {
		name string
		uid  string
		want error
	}{
		{name: "empty", uid: "", want: ErrEmptyUserID},
		{name: "unicode whitespace", uid: "\u2003", want: ErrEmptyUserID},
		{name: "invalid utf8", uid: invalidUTF8, want: ErrInvalidUserID},
		{name: "control", uid: "user\x00name", want: ErrInvalidUserID},
		{name: "too long", uid: strings.Repeat("a", 257), want: ErrInvalidUserID},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := client.FetchPage(
				context.Background(),
				test.uid,
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("transport calls = %d", calls.Load())
	}
}

func TestCustomClientJarIsRemovedAndRedirectIsNotFollowed(t *testing.T) {
	var initialCalls atomic.Int64
	var followCalls atomic.Int64
	var serverURL string
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v0/users/uid/collections":
			initialCalls.Add(1)
			if request.Header.Get("Cookie") != "" {
				t.Errorf("Cookie = %q", request.Header.Get("Cookie"))
			}
			http.Redirect(writer, request, serverURL+"/follow?secret=location-marker", http.StatusFound)
		case "/follow":
			followCalls.Add(1)
			writeJSON(t, writer, fixturePage(0, 1, 0))
		default:
			http.NotFound(writer, request)
		}
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	serverURL = server.URL
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(server.URL)
	jar.SetCookies(parsed, []*http.Cookie{{Name: "secret", Value: "cookie-marker"}})
	original := &http.Client{
		Transport: loopbackOnlyTransport{base: server.Client().Transport},
		Jar:       jar,
		Timeout:   5 * time.Second,
	}
	client := NewClient(
		"collection-tests/0.1",
		WithEndpoint(server.URL),
		WithHTTPClient(original),
		WithRateLimit(1_000_000, 10),
		WithMaxRetries(0),
	)

	_, err = client.FetchPage(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
		1,
		0,
	)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusFound {
		t.Fatalf("error = %T %v", err, err)
	}
	if initialCalls.Load() != 1 || followCalls.Load() != 0 {
		t.Fatalf("initial calls = %d, follow calls = %d", initialCalls.Load(), followCalls.Load())
	}
	if original.Jar != jar || original.Timeout != 5*time.Second {
		t.Fatal("caller-owned client was mutated")
	}
	if strings.Contains(err.Error(), "location-marker") {
		t.Fatalf("Location leaked: %v", err)
	}
}

func TestAttemptTimeoutIsFreshAndCustomClientTimeoutIsIgnored(t *testing.T) {
	var deadlinesMu sync.Mutex
	var deadlines []time.Time
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		deadline, ok := request.Context().Deadline()
		if !ok {
			t.Fatal("attempt has no deadline")
		}
		deadlinesMu.Lock()
		deadlines = append(deadlines, deadline)
		deadlinesMu.Unlock()
		<-request.Context().Done()
		return nil, request.Context().Err()
	})
	original := &http.Client{Transport: transport, Timeout: time.Hour}
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(original),
		WithRequestTimeout(15*time.Millisecond),
		WithMaxRetries(1),
		WithRetryInterval(time.Nanosecond),
		WithRateLimit(1_000_000, 10),
	)
	client.sleep = func(context.Context, time.Duration) error { return nil }

	_, err := client.FetchPage(
		context.Background(),
		"uid",
		SubjectTypeAnime,
		CollectionTypeDone,
		1,
		0,
	)
	var retryErr *RetryError
	if !errors.As(err, &retryErr) || retryErr.Attempts != 2 ||
		!errors.Is(err, ErrTimeout) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %T %v", err, err)
	}
	deadlinesMu.Lock()
	defer deadlinesMu.Unlock()
	if len(deadlines) != 2 || !deadlines[1].After(deadlines[0]) {
		t.Fatalf("deadlines = %v", deadlines)
	}
	if original.Timeout != time.Hour {
		t.Fatal("caller client timeout was mutated")
	}
}

func TestParentDeadlinePrecedesAttemptTimeout(t *testing.T) {
	var calls atomic.Int64
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			calls.Add(1)
			<-request.Context().Done()
			return nil, request.Context().Err()
		})}),
		WithRequestTimeout(time.Second),
		WithMaxRetries(3),
		WithRateLimit(1_000_000, 10),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	_, err := client.FetchPage(ctx, "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	var networkErr *NetworkError
	var retryErr *RetryError
	if !errors.As(err, &networkErr) || !networkErr.Timeout ||
		!errors.Is(err, ErrTimeout) ||
		!errors.Is(err, context.DeadlineExceeded) ||
		errors.As(err, &retryErr) {
		t.Fatalf("error = %T %v", err, err)
	}
	if calls.Load() != 1 {
		t.Fatalf("transport calls = %d", calls.Load())
	}
}

type observedBody struct {
	reader io.Reader
	read   atomic.Int64
	closed atomic.Bool
}

func (body *observedBody) Read(buffer []byte) (int, error) {
	n, err := body.reader.Read(buffer)
	body.read.Add(int64(n))
	return n, err
}

func (body *observedBody) Close() error {
	body.closed.Store(true)
	return nil
}

func TestSuccessAndErrorBodyLimitsAndClosure(t *testing.T) {
	t.Run("oversized success", func(t *testing.T) {
		body := &observedBody{reader: bytes.NewReader(bytes.Repeat([]byte("x"), maxSuccessBodyBytes+32))}
		client := clientWithResponseBody(body, http.StatusOK, nil)
		_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
		if !errors.Is(err, ErrResponseTooLarge) {
			t.Fatalf("error = %T %v", err, err)
		}
		if !body.closed.Load() || body.read.Load() != maxSuccessBodyBytes+1 {
			t.Fatalf("closed = %v, bytes read = %d", body.closed.Load(), body.read.Load())
		}
	})

	t.Run("oversized retryable error", func(t *testing.T) {
		body := &observedBody{reader: bytes.NewReader(bytes.Repeat([]byte("secret"), maxErrorBodyBytes))}
		client := clientWithResponseBody(body, http.StatusServiceUnavailable, nil)
		_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
		if !errors.Is(err, ErrRetryExhausted) ||
			!errors.Is(err, ErrHTTPStatus) ||
			!errors.Is(err, ErrServerError) {
			t.Fatalf("error = %T %v", err, err)
		}
		if !body.closed.Load() || body.read.Load() != maxErrorBodyBytes+1 {
			t.Fatalf("closed = %v, bytes read = %d", body.closed.Load(), body.read.Load())
		}
	})
}

func clientWithResponseBody(body io.ReadCloser, status int, header http.Header) *Client {
	if header == nil {
		header = make(http.Header)
	}
	return NewClient(
		"collection-tests/0.1",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Header:     header,
				Body:       body,
			}, nil
		})}),
		WithRateLimit(1_000_000, 10),
		WithMaxRetries(0),
	)
}

type failingReader struct {
	marker string
}

func (reader failingReader) Read([]byte) (int, error) {
	return 0, errors.New(reader.marker)
}

func TestUnreadableSuccessBodyIsSanitizedDecodeError(t *testing.T) {
	body := &observedBody{reader: failingReader{marker: "raw-body-reader-secret"}}
	client := clientWithResponseBody(body, http.StatusOK, nil)
	_, err := client.FetchPage(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone, 1, 0)
	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) || !errors.Is(err, ErrDecode) {
		t.Fatalf("error = %T %v", err, err)
	}
	if strings.Contains(err.Error(), "raw-body-reader-secret") || !body.closed.Load() {
		t.Fatalf("error leaked or body open: %v, closed=%v", err, body.closed.Load())
	}
}
