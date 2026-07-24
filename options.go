package collection

import (
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Option configures a Client before it becomes available to callers.
type Option func(*Client)

// WithHTTPClient supplies a custom HTTP client.
//
// The supplied value is shallow-cloned. The package-owned clone has no cookie
// jar, has no Client.Timeout, and refuses redirects at the first response.
// WithRequestTimeout remains authoritative for every attempt regardless of
// option order. A nil client preserves the default for compatibility.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = cloneHTTPClient(client)
		}
	}
}

// WithConcurrencyLimit sets the maximum number of in-flight requests shared by
// all operations on one Client. Non-positive values preserve the default.
func WithConcurrencyLimit(limit int) Option {
	return func(c *Client) {
		if limit > 0 {
			c.concurrencyLimit = limit
		}
	}
}

// WithRequestTimeout sets the timeout for each individual HTTP attempt.
// Non-positive values preserve the default.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.requestTimeout = timeout
		}
	}
}

// WithMaxRetries sets retries after the initial attempt. A negative value
// preserves the default; zero disables retries.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		if maxRetries >= 0 {
			c.maxRetries = maxRetries
		}
	}
}

// WithRetryInterval sets the base exponential-backoff interval. Non-positive
// values preserve the default.
func WithRetryInterval(interval time.Duration) Option {
	return func(c *Client) {
		if interval > 0 {
			c.retryInterval = interval
		}
	}
}

// WithEndpoint replaces the API root. It accepts an absolute HTTPS root, or an
// HTTP loopback root for local tests. Invalid values poison the Client with
// ErrInvalidConfiguration.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		parsed, err := parseEndpoint(endpoint)
		if err != nil {
			c.recordConfigurationError()
			return
		}
		c.endpoint = parsed
	}
}

// WithRateLimit sets the shared token-bucket rate and burst. Both values must
// be finite and positive.
func WithRateLimit(requestsPerSecond float64, burst int) Option {
	return func(c *Client) {
		if math.IsNaN(requestsPerSecond) || math.IsInf(requestsPerSecond, 0) ||
			requestsPerSecond <= 0 || burst <= 0 {
			c.recordConfigurationError()
			return
		}
		c.requestsPerSec = requestsPerSecond
		c.rateBurst = burst
	}
}

// WithMaxRetryDelay caps every local and Retry-After-derived retry wait.
func WithMaxRetryDelay(delay time.Duration) Option {
	return func(c *Client) {
		if delay <= 0 {
			c.recordConfigurationError()
			return
		}
		c.maxRetryDelay = delay
	}
}

func cloneHTTPClient(source *http.Client) *http.Client {
	cloned := *source
	cloned.Jar = nil
	cloned.Timeout = 0
	cloned.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &cloned
}

func parseEndpoint(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil ||
		!parsed.IsAbs() ||
		parsed.Opaque != "" ||
		parsed.Host == "" ||
		parsed.User != nil ||
		parsed.RawQuery != "" ||
		parsed.Fragment != "" ||
		(parsed.Path != "" && parsed.Path != "/") ||
		parsed.RawPath != "" {
		return nil, ErrInvalidConfiguration
	}

	switch strings.ToLower(parsed.Scheme) {
	case "https":
	case "http":
		if !isLoopbackHost(parsed.Hostname()) {
			return nil, ErrInvalidConfiguration
		}
	default:
		return nil, ErrInvalidConfiguration
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Path = ""
	parsed.RawPath = ""
	return parsed, nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
