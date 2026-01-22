package collection

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidUserID 用户 ID 不存在或无效
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrUnauthorized 访问令牌无效或已过期
	ErrUnauthorized = errors.New("unauthorized: invalid or expired access token")

	// ErrForbidden 没有权限访问该资源
	ErrForbidden = errors.New("forbidden: access denied")

	// ErrRateLimited 请求过于频繁，触发了速率限制
	ErrRateLimited = errors.New("rate limited: too many requests")

	// ErrServerError 服务器内部错误
	ErrServerError = errors.New("server error")

	// ErrEmptyUserID 用户 ID 为空
	ErrEmptyUserID = errors.New("user id cannot be empty")
)

// NetworkError 网络错误，包装原始错误
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network failed: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// HTTPError 非预期的 HTTP 状态码错误
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("unexpected http status %d: %s", e.StatusCode, e.Body)
}