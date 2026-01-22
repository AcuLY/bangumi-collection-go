package collection

import (
	"net/http"
	"time"
)

// Option 定义客户端配置选项
type Option func(*Client)

// WithHTTPClient 设置自定义 HTTP 客户端
//
// 注意：使用自定义 HTTP 客户端时，WithRequestTimeout 选项将不生效，
// 需要在自定义客户端中自行配置超时
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithAccessToken 设置 Bangumi API 访问令牌
func WithAccessToken(token string) Option {
	return func(c *Client) {
		c.accessToken = token
	}
}

// WithConcurrencyLimit 设置并发请求数限制
func WithConcurrencyLimit(limit int) Option {
	return func(c *Client) {
		if limit > 0 {
			c.concurrencyLimit = limit
		}
	}
}

// WithRequestTimeout 设置单次请求超时时间
//
// 默认值为 30 秒，仅在未使用 WithHTTPClient 时生效
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.requestTimeout = timeout
		}
	}
}

// WithMaxRetries 设置请求失败时的最大重试次数
//
// 默认值为 3 次，仅对可重试的错误（网络错误等）进行重试
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		if maxRetries >= 0 {
			c.maxRetries = maxRetries
		}
	}
}

// WithRetryInterval 设置重试间隔时间
//
// 默认值为 1 秒，实际等待时间会根据重试次数指数增长
func WithRetryInterval(interval time.Duration) Option {
	return func(c *Client) {
		if interval > 0 {
			c.retryInterval = interval
		}
	}
}
