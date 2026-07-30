package collection

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrRateLimited           = errors.New("rate limited")
	ErrServerError           = errors.New("server error")
	ErrEmptyUserID           = errors.New("user id cannot be empty")
	ErrNilContext            = errors.New("nil context")
	ErrNoCollectionTypes     = errors.New("no collection types")
	ErrInvalidSubjectType    = errors.New("invalid subject type")
	ErrInvalidCollectionType = errors.New("invalid collection type")
	ErrInvalidConfiguration  = errors.New("invalid client configuration")
	ErrNotFound              = errors.New("not found")
	ErrHTTPStatus            = errors.New("unexpected http status")
	ErrTransport             = errors.New("transport failure")
	ErrTimeout               = errors.New("request timeout")
	ErrCanceled              = errors.New("request canceled")
	ErrDecode                = errors.New("response decode failure")
	ErrProtocol              = errors.New("response protocol violation")
	ErrResponseTooLarge      = errors.New("response too large")
	ErrRetryExhausted        = errors.New("retry exhausted")
)

// HTTPError describes a non-200 response without retaining response content.
// Body remains for source compatibility and is always empty on returned errors.
type HTTPError struct {
	StatusCode int
	Body       string
	RetryAfter time.Duration
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("collection request failed: http status %d", e.StatusCode)
}

func (e *HTTPError) Is(target error) bool {
	if target == ErrHTTPStatus {
		return true
	}
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return target == ErrUnauthorized
	case http.StatusForbidden:
		return target == ErrForbidden
	case http.StatusNotFound:
		return target == ErrNotFound || target == ErrInvalidUserID
	case http.StatusTooManyRequests:
		return target == ErrRateLimited
	default:
		return e.StatusCode >= 500 && e.StatusCode <= 599 && target == ErrServerError
	}
}

// NetworkError describes a sanitized runtime request failure.
//
// Returned Err values are restricted to context.Canceled,
// context.DeadlineExceeded, or ErrTransport.
type NetworkError struct {
	Err      error
	Timeout  bool
	terminal bool
}

func (e *NetworkError) Error() string {
	switch {
	case e.Timeout:
		return "collection request failed: timeout"
	case errors.Is(e.Err, context.Canceled):
		return "collection request failed: canceled"
	default:
		return "collection request failed: transport"
	}
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

func (e *NetworkError) Is(target error) bool {
	switch {
	case target == ErrTimeout && e.Timeout:
		return true
	case target == ErrCanceled && errors.Is(e.Err, context.Canceled):
		return true
	case target == ErrTransport && !e.Timeout && !errors.Is(e.Err, context.Canceled):
		return true
	default:
		return false
	}
}

// DecodeError indicates that a bounded success body could not be read or was
// not exactly one JSON value. It intentionally exposes no upstream content.
type DecodeError struct {
	kind string
}

func (e *DecodeError) Error() string {
	return "collection response decode failed"
}

func (e *DecodeError) Unwrap() error {
	return ErrDecode
}

func newDecodeError(kind string) *DecodeError {
	return &DecodeError{kind: kind}
}

// ProtocolError indicates that decoded data violated the collection contract.
// It intentionally exposes no upstream values.
type ProtocolError struct{}

func (e *ProtocolError) Error() string {
	return "collection response protocol violation"
}

func (e *ProtocolError) Unwrap() error {
	return ErrProtocol
}

func newProtocolError() *ProtocolError {
	return &ProtocolError{}
}

type responseTooLargeError struct{}

func (e *responseTooLargeError) Error() string {
	return "collection response too large"
}

func (e *responseTooLargeError) Unwrap() error {
	return ErrResponseTooLarge
}

// RetryError reports exhaustion while preserving the last sanitized error.
type RetryError struct {
	Attempts int
	Err      error
}

func (e *RetryError) Error() string {
	return fmt.Sprintf("collection request retry exhausted after %d attempts", e.Attempts)
}

func (e *RetryError) Unwrap() error {
	return e.Err
}

func (e *RetryError) Is(target error) bool {
	return target == ErrRetryExhausted
}

func networkErrorForContext(err error) *NetworkError {
	if errors.Is(err, context.DeadlineExceeded) {
		return &NetworkError{Err: context.DeadlineExceeded, Timeout: true, terminal: true}
	}
	return &NetworkError{Err: context.Canceled, Timeout: false, terminal: true}
}

func transportError() *NetworkError {
	return &NetworkError{Err: ErrTransport}
}
