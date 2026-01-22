// Package collection 提供获取 Bangumi（番组计划）用户收藏列表的 Go 客户端。
//
// 本包支持并发获取用户的动画、书籍、游戏、音乐等条目的收藏数据，
// 包括想看、在看、看过、搁置、抛弃等收藏状态。
//
// 基本用法:
//
//	client := collection.NewClient("AcuL/my-private-project")
//
//	subjects, err := client.Fetch(
// 		ctx, 
// 		"username", 
// 		collection.SubjectTypeAnime,
//	    collection.CollectionTypeDone,
//	    collection.CollectionTypeDoing,
//	)
//
// 支持用选项函数附加 Bangumi 的访问令牌、自定义 HTTP 客户端、并发限制、超时重试等
package collection

import (
	"net/http"
	"time"
)

const (
	defaultConcurrencyLimit = 10
	defaultRequestTimeout   = 30 * time.Second
	defaultMaxRetries       = 3
	defaultRetryInterval    = time.Second
)

// Client 抓取 Bangumi 收藏列表客户端
type Client struct {
	httpClient       *http.Client
	userAgent        string
	accessToken      string
	concurrencyLimit int
	requestTimeout   time.Duration
	maxRetries       int
	retryInterval    time.Duration
}

// NewClient 创建新的 Bangumi 收藏抓取客户端
//
// userAgent 必填，
// 参考 https://github.com/bangumi/api/blob/master/docs-raw/user%20agent.md
func NewClient(userAgent string, options ...Option) *Client {
	c := &Client{
		userAgent:        userAgent,
		concurrencyLimit: defaultConcurrencyLimit,
		requestTimeout:   defaultRequestTimeout,
		maxRetries:       defaultMaxRetries,
		retryInterval:    defaultRetryInterval,
	}

	for _, opt := range options {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.requestTimeout,
		}
	}

	return c
}
