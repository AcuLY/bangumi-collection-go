package collection

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *Client) doRequest(ctx context.Context, params fetchParams) (*PageResult, error) {
	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, networkErrorForContext(err)
		}

		page, err := c.doRequestAttempt(ctx, params)
		if err == nil {
			return page, nil
		}
		if contextErr := ctx.Err(); contextErr != nil {
			return nil, networkErrorForContext(contextErr)
		}
		if !retryable(err) {
			return nil, err
		}
		if attempt > c.maxRetries || attempt == maxInt() {
			return nil, &RetryError{Attempts: attempt, Err: err}
		}

		retryAfter := time.Duration(0)
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			retryAfter = httpErr.RetryAfter
		}
		delay := c.retryDelay(attempt, retryAfter)
		if err := c.sleep(ctx, delay); err != nil {
			if contextErr := ctx.Err(); contextErr != nil {
				return nil, networkErrorForContext(contextErr)
			}
			return nil, networkErrorForContext(context.Canceled)
		}
	}
}

func retryable(err error) bool {
	var networkErr *NetworkError
	if errors.As(err, &networkErr) {
		return !networkErr.terminal &&
			(errors.Is(err, ErrTransport) || errors.Is(err, ErrTimeout))
	}
	var httpErr *HTTPError
	return errors.As(err, &httpErr) &&
		(httpErr.StatusCode == http.StatusTooManyRequests ||
			(httpErr.StatusCode >= 500 && httpErr.StatusCode <= 599))
}

func (c *Client) retryDelay(retryNumber int, retryAfter time.Duration) time.Duration {
	capDelay := c.retryInterval
	for exponent := 1; exponent < retryNumber && capDelay < c.maxRetryDelay; exponent++ {
		if capDelay > c.maxRetryDelay/2 {
			capDelay = c.maxRetryDelay
			break
		}
		capDelay *= 2
	}
	if capDelay > c.maxRetryDelay {
		capDelay = c.maxRetryDelay
	}

	random := c.randomFloat()
	if math.IsNaN(random) || random < 0 {
		random = 0
	}
	if random > 1 {
		random = 1
	}
	jitter := capDelay
	if random < 1 {
		jitter = time.Duration(random * float64(capDelay))
		if jitter < 0 || jitter > capDelay {
			jitter = capDelay
		}
	}
	if retryAfter > c.maxRetryDelay {
		retryAfter = c.maxRetryDelay
	}
	if retryAfter < 0 {
		retryAfter = 0
	}
	if retryAfter > jitter {
		return retryAfter
	}
	return jitter
}

func (c *Client) parseRetryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	deltaSeconds := value != ""
	for _, character := range value {
		if character < '0' || character > '9' {
			deltaSeconds = false
			break
		}
	}
	if deltaSeconds {
		digits := strings.TrimLeft(value, "0")
		if digits == "" {
			return 0
		}
		maxSeconds := uint64(c.maxRetryDelay / time.Second)
		maxText := strconv.FormatUint(maxSeconds, 10)
		if len(digits) > len(maxText) ||
			(len(digits) == len(maxText) && digits > maxText) {
			return c.maxRetryDelay
		}
		seconds, err := strconv.ParseUint(digits, 10, 64)
		if err != nil {
			return c.maxRetryDelay
		}
		delay := time.Duration(seconds) * time.Second
		if delay > c.maxRetryDelay {
			return c.maxRetryDelay
		}
		return delay
	}

	parsed, err := http.ParseTime(value)
	if err != nil {
		return 0
	}
	delay := parsed.Sub(c.now())
	if delay <= 0 {
		return 0
	}
	if delay > c.maxRetryDelay {
		return c.maxRetryDelay
	}
	return delay
}
