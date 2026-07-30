package collection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type loopbackOnlyTransport struct {
	base http.RoundTripper
}

func (transport loopbackOnlyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	host := request.URL.Hostname()
	ip := net.ParseIP(host)
	if !strings.EqualFold(host, "localhost") && (ip == nil || !ip.IsLoopback()) {
		return nil, fmt.Errorf("test blocked non-loopback request")
	}
	base := transport.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(request)
}

func newLoopbackServerClient(
	t *testing.T,
	handler http.Handler,
	options ...Option,
) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	baseOptions := []Option{
		WithEndpoint(server.URL),
		WithHTTPClient(&http.Client{
			Transport: loopbackOnlyTransport{base: server.Client().Transport},
		}),
		WithRateLimit(1_000_000, 100_000),
		WithMaxRetries(0),
	}
	baseOptions = append(baseOptions, options...)
	return NewClient("collection-tests/0.1", baseOptions...), server
}

func fixtureCollection(subjectID int, subjectType SubjectType, collectionType CollectionType) map[string]any {
	return map[string]any{
		"updated_at":   "2026-07-24T12:34:56Z",
		"comment":      "fixture comment",
		"tags":         []string{"one", "two"},
		"subject_id":   subjectID,
		"vol_status":   2,
		"ep_status":    7,
		"subject_type": int(subjectType),
		"type":         int(collectionType),
		"rate":         8,
		"private":      false,
		"subject": map[string]any{
			"id":      subjectID,
			"type":    int(subjectType),
			"name":    fmt.Sprintf("subject-%d", subjectID),
			"name_cn": fmt.Sprintf("条目-%d", subjectID),
		},
	}
}

func fixturePage(
	total int,
	limit int,
	offset int,
	items ...map[string]any,
) []byte {
	if items == nil {
		items = make([]map[string]any, 0)
	}
	body, err := json.Marshal(map[string]any{
		"data":   items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		panic(err)
	}
	return body
}

func writeJSON(t *testing.T, writer http.ResponseWriter, body []byte) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(body); err != nil {
		t.Errorf("write response: %v", err)
	}
}

func TestClientDefaults(t *testing.T) {
	client := NewClient("collection-tests/0.1")
	if client.endpoint.String() != defaultEndpoint {
		t.Fatalf("endpoint = %q", client.endpoint.String())
	}
	if client.concurrencyLimit != defaultConcurrencyLimit ||
		client.requestTimeout != defaultRequestTimeout ||
		client.maxRetries != defaultMaxRetries ||
		client.retryInterval != defaultRetryInterval ||
		client.maxRetryDelay != defaultMaxRetryDelay ||
		client.requestsPerSec != defaultRequestsPerSec ||
		client.rateBurst != defaultRateBurst {
		t.Fatalf("unexpected defaults: %#v", client)
	}
	if client.httpClient == http.DefaultClient {
		t.Fatal("default HTTP client was not cloned")
	}
	if client.httpClient.Timeout != 0 || client.httpClient.Jar != nil {
		t.Fatal("package HTTP client retained timeout or jar")
	}
	if client.configErr != nil {
		t.Fatalf("config error = %v", client.configErr)
	}
}

func TestInvalidUserAgentFailsBeforeTransport(t *testing.T) {
	invalidUTF8 := string([]byte{0xff})
	tests := []string{
		"",
		" \u2003 ",
		"bad\nagent",
		invalidUTF8,
		strings.Repeat("a", 257),
	}
	for _, userAgent := range tests {
		t.Run(fmt.Sprintf("%x", []byte(userAgent)), func(t *testing.T) {
			var calls atomic.Int64
			client := NewClient(
				userAgent,
				WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					calls.Add(1)
					return nil, errors.New("must not run")
				})}),
			)
			_, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				1,
				0,
			)
			if !errors.Is(err, ErrInvalidConfiguration) {
				t.Fatalf("error = %v", err)
			}
			if calls.Load() != 0 {
				t.Fatalf("transport calls = %d", calls.Load())
			}
		})
	}
}

func TestNilOptionPoisonsClientAndCannotBeCleared(t *testing.T) {
	client := NewClient(
		"collection-tests/0.1",
		nil,
		WithRateLimit(10, 2),
		WithMaxRetryDelay(time.Second),
	)
	if !errors.Is(client.configErr, ErrInvalidConfiguration) {
		t.Fatalf("config error = %v", client.configErr)
	}
	_, err := client.Fetch(context.Background(), "uid", SubjectTypeAnime, CollectionTypeDone)
	if !errors.Is(err, ErrInvalidConfiguration) {
		t.Fatalf("operation error = %v", err)
	}
}

func TestAcceptedUserAgentIsSentByteForByte(t *testing.T) {
	const userAgent = "\u2003AcuL/collection-client\u2003"
	client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("User-Agent"); got != userAgent {
			t.Errorf("User-Agent = %q", got)
		}
		writeJSON(t, writer, fixturePage(0, 1, 0))
	}))
	client.userAgent = userAgent

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
}
