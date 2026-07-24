package collection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	maxSuccessBodyBytes = 16 << 20
	maxErrorBodyBytes   = 64 << 10
	maxPageTotal        = 1_000_000
)

func normalizeUserID(value string) (string, error) {
	if !utf8.ValidString(value) {
		return "", ErrInvalidUserID
	}
	normalized := strings.TrimFunc(value, unicode.IsSpace)
	if normalized == "" {
		return "", ErrEmptyUserID
	}
	if len(normalized) > 256 {
		return "", ErrInvalidUserID
	}
	for _, r := range normalized {
		if unicode.IsControl(r) {
			return "", ErrInvalidUserID
		}
	}
	return normalized, nil
}

func (c *Client) buildRequest(ctx context.Context, params fetchParams) (*http.Request, error) {
	root := *c.endpoint
	root.Path = "/v0/users/" + params.UserID + "/collections"
	root.RawPath = "/v0/users/" + escapePathSegment(params.UserID) + "/collections"

	query := make(url.Values, 4)
	query.Set("subject_type", strconv.Itoa(int(params.SubjectType)))
	query.Set("type", strconv.Itoa(int(params.CollectionType)))
	query.Set("limit", strconv.Itoa(params.Limit))
	query.Set("offset", strconv.Itoa(params.Offset))
	root.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, root.String(), nil)
	if err != nil {
		return nil, ErrInvalidConfiguration
	}
	request.Header.Set("User-Agent", c.userAgent)
	request.Header.Set("Accept", "application/json")
	return request, nil
}

func escapePathSegment(value string) string {
	switch value {
	case ".":
		return "%2E"
	case "..":
		return "%2E%2E"
	default:
		return url.PathEscape(value)
	}
}

func (c *Client) doRequestAttempt(ctx context.Context, params fetchParams) (*PageResult, error) {
	release, err := c.acquireAttempt(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	attemptCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	request, err := c.buildRequest(attemptCtx, params)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		if classified := classifyContextFailure(ctx, attemptCtx); classified != nil {
			return nil, classified
		}
		return nil, transportError()
	}

	if response.StatusCode != http.StatusOK {
		readErr := discardAndClose(response.Body, maxErrorBodyBytes+1)
		if classified := classifyContextFailure(ctx, attemptCtx); classified != nil {
			return nil, classified
		}
		_ = readErr
		retryAfter := time.Duration(0)
		if response.StatusCode == http.StatusTooManyRequests ||
			(response.StatusCode >= 500 && response.StatusCode <= 599) {
			retryAfter = c.parseRetryAfter(response.Header.Get("Retry-After"))
		}
		return nil, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       "",
			RetryAfter: retryAfter,
		}
	}

	body, readErr := readAndClose(response.Body, maxSuccessBodyBytes+1)
	if classified := classifyContextFailure(ctx, attemptCtx); classified != nil {
		return nil, classified
	}
	if readErr != nil {
		return nil, newDecodeError("read")
	}
	if len(body) > maxSuccessBodyBytes {
		return nil, &responseTooLargeError{}
	}
	page, err := decodePage(body, params)
	if classified := classifyContextFailure(ctx, attemptCtx); classified != nil {
		return nil, classified
	}
	return page, err
}

func classifyContextFailure(parent, attempt context.Context) error {
	if err := parent.Err(); err != nil {
		return networkErrorForContext(err)
	}
	if err := attempt.Err(); err != nil {
		return &NetworkError{Err: context.DeadlineExceeded, Timeout: true}
	}
	return nil
}

func readAndClose(body io.ReadCloser, limit int64) ([]byte, error) {
	if body == nil {
		return nil, io.ErrUnexpectedEOF
	}
	data, readErr := io.ReadAll(io.LimitReader(body, limit))
	closeErr := body.Close()
	if readErr != nil {
		return data, readErr
	}
	return data, closeErr
}

func discardAndClose(body io.ReadCloser, limit int64) error {
	if body == nil {
		return nil
	}
	_, readErr := io.Copy(io.Discard, io.LimitReader(body, limit))
	closeErr := body.Close()
	if readErr != nil {
		return readErr
	}
	return closeErr
}

func decodePage(body []byte, params fetchParams) (*PageResult, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	var wire wirePage
	if err := decoder.Decode(&wire); err != nil {
		return nil, newDecodeError("json")
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, newDecodeError("trailing")
	}

	if wire.Data == nil || wire.Total == nil || wire.Limit == nil || wire.Offset == nil {
		return nil, newProtocolError()
	}
	if *wire.Total < 0 || *wire.Total > maxPageTotal ||
		*wire.Limit != int64(params.Limit) ||
		*wire.Offset != int64(params.Offset) ||
		len(*wire.Data) > params.Limit {
		return nil, newProtocolError()
	}

	data := make([]*Subject, 0, len(*wire.Data))
	for _, item := range *wire.Data {
		subject, err := item.toSubject(params.SubjectType, params.CollectionType)
		if err != nil {
			return nil, err
		}
		data = append(data, subject)
	}
	if data == nil {
		data = make([]*Subject, 0)
	}

	return &PageResult{
		Data:   data,
		Total:  int(*wire.Total),
		Limit:  int(*wire.Limit),
		Offset: int(*wire.Offset),
	}, nil
}
