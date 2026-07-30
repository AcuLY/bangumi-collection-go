package collection

import (
	"context"
	"errors"
)

func (c *Client) acquireAttempt(ctx context.Context) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, networkErrorForContext(err)
	}

	if err := c.requestLimiter.Wait(ctx); err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return nil, networkErrorForContext(contextErr)
		}
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			return nil, networkErrorForContext(context.DeadlineExceeded)
		}
		if errors.Is(err, context.Canceled) {
			return nil, networkErrorForContext(context.Canceled)
		}
		return nil, transportError()
	}

	select {
	case <-ctx.Done():
		return nil, networkErrorForContext(ctx.Err())
	case c.inflight <- struct{}{}:
	}
	if err := ctx.Err(); err != nil {
		<-c.inflight
		return nil, networkErrorForContext(err)
	}

	released := false
	return func() {
		if released {
			return
		}
		released = true
		<-c.inflight
	}, nil
}
