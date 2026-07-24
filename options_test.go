package collection

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"testing"
	"time"
)

func TestExistingOptionFallbacksAndCallShapes(t *testing.T) {
	client := NewClient(
		"collection-tests/0.1",
		WithHTTPClient(nil),
		WithConcurrencyLimit(0),
		WithRequestTimeout(0),
		WithMaxRetries(-1),
		WithRetryInterval(0),
	)
	if client.concurrencyLimit != defaultConcurrencyLimit ||
		client.requestTimeout != defaultRequestTimeout ||
		client.maxRetries != defaultMaxRetries ||
		client.retryInterval != defaultRetryInterval {
		t.Fatalf("fallback defaults changed: %#v", client)
	}

	var _ func(string, ...Option) *Client = NewClient
	var _ func(*http.Client) Option = WithHTTPClient
	var _ func(int) Option = WithConcurrencyLimit
	var _ func(time.Duration) Option = WithRequestTimeout
	var _ func(int) Option = WithMaxRetries
	var _ func(time.Duration) Option = WithRetryInterval
	var _ func(string) Option = WithEndpoint
	var _ func(float64, int) Option = WithRateLimit
	var _ func(time.Duration) Option = WithMaxRetryDelay
	var _ func(*Client, context.Context, string, SubjectType, ...CollectionType) ([]*Subject, error) = (*Client).Fetch
}

func TestInvalidNewOptionsPoisonClient(t *testing.T) {
	tests := []struct {
		name   string
		option Option
	}{
		{name: "nan qps", option: WithRateLimit(math.NaN(), 1)},
		{name: "positive infinity", option: WithRateLimit(math.Inf(1), 1)},
		{name: "negative infinity", option: WithRateLimit(math.Inf(-1), 1)},
		{name: "zero qps", option: WithRateLimit(0, 1)},
		{name: "negative qps", option: WithRateLimit(-1, 1)},
		{name: "zero burst", option: WithRateLimit(1, 0)},
		{name: "negative burst", option: WithRateLimit(1, -1)},
		{name: "zero max delay", option: WithMaxRetryDelay(0)},
		{name: "negative max delay", option: WithMaxRetryDelay(-time.Second)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(
				"collection-tests/0.1",
				test.option,
				WithRateLimit(50, 5),
				WithMaxRetryDelay(time.Second),
			)
			if !errors.Is(client.configErr, ErrInvalidConfiguration) {
				t.Fatalf("config error = %v", client.configErr)
			}
		})
	}
}

func TestEndpointValidation(t *testing.T) {
	valid := []string{
		"https://api.bgm.tv",
		"https://api.bgm.tv/",
		"https://example.test:8443",
		"http://127.0.0.1:8080",
		"http://[::1]:8080",
		"http://localhost:8080",
	}
	for _, endpoint := range valid {
		t.Run("valid "+endpoint, func(t *testing.T) {
			client := NewClient("collection-tests/0.1", WithEndpoint(endpoint))
			if client.configErr != nil {
				t.Fatalf("config error = %v", client.configErr)
			}
		})
	}

	invalid := []string{
		"",
		"/relative",
		"https://",
		"https://user:password@example.test",
		"https://example.test/v0",
		"https://example.test/?query=1",
		"https://example.test/#fragment",
		"http://example.test",
		"ftp://example.test",
		"https:opaque",
	}
	for _, endpoint := range invalid {
		t.Run("invalid "+endpoint, func(t *testing.T) {
			client := NewClient(
				"collection-tests/0.1",
				WithEndpoint(endpoint),
				WithEndpoint("https://api.bgm.tv"),
			)
			if !errors.Is(client.configErr, ErrInvalidConfiguration) {
				t.Fatalf("config error = %v", client.configErr)
			}
		})
	}
}

func TestHTTPClientIsClonedAndTimeoutOptionOrderIsIndependent(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	redirect := func(*http.Request, []*http.Request) error { return nil }
	original := &http.Client{
		Transport:     roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("unused") }),
		CheckRedirect: redirect,
		Jar:           jar,
		Timeout:       11 * time.Second,
	}

	tests := []struct {
		name    string
		options []Option
	}{
		{
			name: "client then timeout",
			options: []Option{
				WithHTTPClient(original),
				WithRequestTimeout(17 * time.Millisecond),
			},
		},
		{
			name: "timeout then client",
			options: []Option{
				WithRequestTimeout(17 * time.Millisecond),
				WithHTTPClient(original),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient("collection-tests/0.1", test.options...)
			if client.httpClient == original {
				t.Fatal("custom client was not cloned")
			}
			if client.httpClient.Timeout != 0 || client.httpClient.Jar != nil {
				t.Fatal("clone retained timeout or jar")
			}
			if client.httpClient.CheckRedirect == nil {
				t.Fatal("clone has no redirect refusal")
			}
			if client.requestTimeout != 17*time.Millisecond {
				t.Fatalf("request timeout = %s", client.requestTimeout)
			}
		})
	}

	if original.Timeout != 11*time.Second || original.Jar != jar ||
		reflect.ValueOf(original.CheckRedirect).Pointer() != reflect.ValueOf(redirect).Pointer() {
		t.Fatal("caller-owned HTTP client was mutated")
	}
}
