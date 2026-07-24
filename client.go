// Package collection provides an anonymous, read-only client for Bangumi
// public user collections.
//
// A Client is safe for concurrent use after construction. Fetch retrieves all
// pages and returns a canonical order; FetchPage preserves one upstream page's
// order. Requests never send Authorization or Cookie headers.
//
// Bangumi collection records require tags but may omit comment and the nested
// subject projection. An omitted comment becomes the empty string; an omitted
// subject keeps ID equal to SubjectID and leaves Name and NameCn empty. Present
// optional fields must still contain a valid non-null value.
//
// Subject types follow the official mapping: Book 1, Anime 2, Music 3, Game 4,
// and Real 6. Music and Game intentionally correct the reversed names in the
// untagged prototype before the first public version.
//
// The first public contract retains NewClient, Fetch, FetchPage, collection
// enum values, and existing non-authentication options. It intentionally
// removes WithAccessToken, rejects an empty Fetch collection-type list, extends
// Subject to the complete collection DTO, and keeps HTTPError.Body empty. New
// options configure a test endpoint, a shared rate limit, and a maximum retry
// delay.
//
// This repository is preparing the first v0.1.0 contract. The package has not
// been tagged or published by this change.
package collection

import (
	"context"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/time/rate"
)

const (
	defaultEndpoint         = "https://api.bgm.tv"
	defaultConcurrencyLimit = 10
	defaultRequestTimeout   = 30 * time.Second
	defaultMaxRetries       = 3
	defaultRetryInterval    = time.Second
	defaultMaxRetryDelay    = 30 * time.Second
	defaultRequestsPerSec   = 3
	defaultRateBurst        = 1
)

// Client retrieves Bangumi public collection data.
//
// Configuration is fixed when NewClient returns. The same Client shares its
// request-rate and in-flight limits across all Fetch and FetchPage calls.
type Client struct {
	httpClient       *http.Client
	endpoint         *url.URL
	userAgent        string
	concurrencyLimit int
	requestTimeout   time.Duration
	maxRetries       int
	retryInterval    time.Duration
	maxRetryDelay    time.Duration
	requestsPerSec   float64
	rateBurst        int
	configErr        error

	requestLimiter *rate.Limiter
	inflight       chan struct{}

	// The hooks are immutable in production. Package tests may replace them
	// before the Client is used to make retry behavior deterministic.
	now         func() time.Time
	sleep       func(context.Context, time.Duration) error
	randomFloat func() float64
}

// NewClient creates an anonymous Bangumi public-collection client.
//
// userAgent is required and must be valid UTF-8, contain no control rune, and
// be between 1 and 256 bytes. Since this constructor preserves its historical
// signature, invalid configuration is retained on the Client and returned as
// ErrInvalidConfiguration by every operation before transport.
func NewClient(userAgent string, options ...Option) *Client {
	endpoint, _ := parseEndpoint(defaultEndpoint)
	c := &Client{
		endpoint:         endpoint,
		userAgent:        userAgent,
		concurrencyLimit: defaultConcurrencyLimit,
		requestTimeout:   defaultRequestTimeout,
		maxRetries:       defaultMaxRetries,
		retryInterval:    defaultRetryInterval,
		maxRetryDelay:    defaultMaxRetryDelay,
		requestsPerSec:   defaultRequestsPerSec,
		rateBurst:        defaultRateBurst,
		now:              time.Now,
		sleep:            sleepContext,
		randomFloat:      rand.Float64,
	}

	if !validUserAgent(userAgent) {
		c.recordConfigurationError()
	}

	for _, option := range options {
		if option == nil {
			c.recordConfigurationError()
			continue
		}
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = cloneHTTPClient(http.DefaultClient)
	}
	c.requestLimiter = rate.NewLimiter(rate.Limit(c.requestsPerSec), c.rateBurst)
	c.inflight = make(chan struct{}, c.concurrencyLimit)

	return c
}

func (c *Client) recordConfigurationError() {
	if c.configErr == nil {
		c.configErr = ErrInvalidConfiguration
	}
}

func validUserAgent(value string) bool {
	if !utf8.ValidString(value) || len(value) == 0 || len(value) > 256 {
		return false
	}
	if strings.TrimFunc(value, unicode.IsSpace) == "" {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func sleepContext(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
