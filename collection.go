package collection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	baseURL  = "https://api.bgm.tv/v0/users/%s/collections"
	pageSize = 50
)

// Fetch 并发获取指定类型的全部收藏
//
// 参考 https://bangumi.github.io/api/#/%E6%94%B6%E8%97%8F/getUserCollectionsByUsername
func (c *Client) Fetch(ctx context.Context, userID string, subjType SubjectType, collTypes ...CollectionType) ([]*Subject, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	var (
		subjects []*Subject
		mu       sync.Mutex
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(c.concurrencyLimit)

	for _, ct := range collTypes {
		g.Go(func() error {
			p := fetchParams{
				UserID:         userID,
				SubjectType:    subjType,
				CollectionType: ct,
			}

			result, err := c.fetchAllPages(ctx, p)
			if err != nil {
				return err
			}

			mu.Lock()
			subjects = append(subjects, result...)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return subjects, nil
}

// FetchPage 获取单页收藏数据，支持自定义 limit 和 offset
//
// limit 范围为 1-50，超出范围会被自动限制
//
// offset 从 0 开始，小于 0 会被设为 0
func (c *Client) FetchPage(ctx context.Context, userID string, subjType SubjectType, collType CollectionType, limit, offset int) (*PageResult, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	p := fetchParams{
		UserID:         userID,
		SubjectType:    subjType,
		CollectionType: collType,
		Limit:          clamp(limit, 1, pageSize),
		Offset:         max(offset, 0),
	}

	r, err := c.doRequest(ctx, p)
	if err != nil {
		return nil, err
	}

	return &PageResult{
		Data:   r.toSubjects(),
		Total:  r.Total,
		Limit:  r.Limit,
		Offset: r.Offset,
	}, nil
}

func validateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrEmptyUserID
	}
	return nil
}

func clamp(v, minVal, maxVal int) int {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

func (c *Client) fetchAllPages(ctx context.Context, p fetchParams) ([]*Subject, error) {
	p.Limit = 1
	firstPage, err := c.doRequest(ctx, p)
	if err != nil {
		return nil, err
	}

	total := firstPage.Total
	if total == 0 {
		return nil, nil
	}

	var (
		subjects = make([]*Subject, 0, total)
		mu       sync.Mutex
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(c.concurrencyLimit)

	for offset := 0; offset < total; offset += pageSize {
		p := fetchParams{
			UserID:         p.UserID,
			SubjectType:    p.SubjectType,
			CollectionType: p.CollectionType,
			Offset:         offset,
			Limit:          pageSize,
		}

		g.Go(func() error {
			r, err := c.doRequest(ctx, p)
			if err != nil {
				return err
			}
			mu.Lock()
			subjects = append(subjects, r.toSubjects()...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return subjects, nil
}

func (c *Client) doRequest(ctx context.Context, p fetchParams) (*result, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// 非首次尝试时等待
		if attempt > 0 {
			waitTime := c.retryInterval * time.Duration(1<<(attempt-1)) // 指数退避
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}
		}

		r, err := c.doRequestOnce(ctx, p)
		if err == nil {
			return r, nil
		}

		// 判断是否可重试
		if !c.isRetryable(err) {
			return nil, err
		}
		lastErr = err
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (c *Client) isRetryable(err error) bool {
	if err == ErrRateLimited || err == ErrServerError {
		return true
	}
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return true
	}
	return false
}

func (c *Client) doRequestOnce(ctx context.Context, p fetchParams) (*result, error) {
	req, err := c.buildRequest(ctx, p)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &NetworkError{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &NetworkError{Err: err}
	}

	// 处理 HTTP 状态码
	if err := c.handleStatusCode(resp.StatusCode, body); err != nil {
		return nil, err
	}

	var r result
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &r, nil
}

func (c *Client) handleStatusCode(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrInvalidUserID
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusTooManyRequests:
		return ErrRateLimited
	default:
		if statusCode >= 500 {
			return ErrServerError
		}
		if statusCode >= 400 {
			return &HTTPError{
				StatusCode: statusCode,
				Body:       string(body),
			}
		}
		return nil
	}
}

func (c *Client) buildRequest(ctx context.Context, p fetchParams) (*http.Request, error) {
	url := fmt.Sprintf(baseURL, p.UserID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("subject_type", strconv.Itoa(int(p.SubjectType)))
	q.Set("type", strconv.Itoa(int(p.CollectionType)))
	q.Set("offset", strconv.Itoa(p.Offset))
	q.Set("limit", strconv.Itoa(p.Limit))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", c.userAgent)
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	return req, nil
}
