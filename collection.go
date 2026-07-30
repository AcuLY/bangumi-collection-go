package collection

import (
	"context"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"
)

const pageSize = 50

// Fetch retrieves all planned pages for the requested collection states.
//
// Repeated states are normalized. Results are deduplicated and sorted by
// (SubjectType, SubjectID, Type), independent of completion order.
func (c *Client) Fetch(
	ctx context.Context,
	userID string,
	subjectType SubjectType,
	collectionTypes ...CollectionType,
) ([]*Subject, error) {
	normalizedID, normalizedTypes, err := c.validateFetchInput(ctx, userID, subjectType, collectionTypes)
	if err != nil {
		return nil, err
	}

	firstPages := make([]locatedPage, 0, len(normalizedTypes))
	jobs := make([]pageJob, 0)
	for _, collectionType := range normalizedTypes {
		params := fetchParams{
			UserID:         normalizedID,
			SubjectType:    subjectType,
			CollectionType: collectionType,
			Limit:          pageSize,
			Offset:         0,
		}
		page, err := c.doRequest(ctx, params)
		if err != nil {
			return nil, err
		}
		if len(page.Data) > page.Total {
			return nil, newProtocolError()
		}
		firstPages = append(firstPages, locatePage(page, collectionType))

		pageCount := 0
		if page.Total > 0 {
			pageCount = (page.Total + pageSize - 1) / pageSize
		}
		if pageCount > 0 && pageCount-1 > maxInt()/pageSize {
			return nil, newProtocolError()
		}
		for pageIndex := 1; pageIndex < pageCount; pageIndex++ {
			offset := pageIndex * pageSize
			if offset < 0 || offset >= page.Total {
				return nil, newProtocolError()
			}
			jobs = append(jobs, pageJob{
				params: fetchParams{
					UserID:         normalizedID,
					SubjectType:    subjectType,
					CollectionType: collectionType,
					Limit:          pageSize,
					Offset:         offset,
				},
			})
		}
	}

	jobPages, err := c.fetchRemainingPages(ctx, jobs)
	if err != nil {
		return nil, err
	}

	all := make([]locatedSubject, 0)
	for _, page := range firstPages {
		all = append(all, page.items...)
	}
	for _, page := range jobPages {
		all = append(all, page.items...)
	}
	return canonicalSubjects(all), nil
}

// FetchPage retrieves one validated page. limit is clamped to 1..50 and a
// negative offset is clamped to zero for compatibility.
func (c *Client) FetchPage(
	ctx context.Context,
	userID string,
	subjectType SubjectType,
	collectionType CollectionType,
	limit int,
	offset int,
) (*PageResult, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if c.configErr != nil {
		return nil, c.configErr
	}
	normalizedID, err := normalizeUserID(userID)
	if err != nil {
		return nil, err
	}
	if !validSubjectType(subjectType) {
		return nil, ErrInvalidSubjectType
	}
	if !validCollectionType(collectionType) {
		return nil, ErrInvalidCollectionType
	}

	return c.doRequest(ctx, fetchParams{
		UserID:         normalizedID,
		SubjectType:    subjectType,
		CollectionType: collectionType,
		Limit:          clamp(limit, 1, pageSize),
		Offset:         max(offset, 0),
	})
}

func (c *Client) validateFetchInput(
	ctx context.Context,
	userID string,
	subjectType SubjectType,
	collectionTypes []CollectionType,
) (string, []CollectionType, error) {
	if ctx == nil {
		return "", nil, ErrNilContext
	}
	if c.configErr != nil {
		return "", nil, c.configErr
	}
	normalizedID, err := normalizeUserID(userID)
	if err != nil {
		return "", nil, err
	}
	if !validSubjectType(subjectType) {
		return "", nil, ErrInvalidSubjectType
	}
	if len(collectionTypes) == 0 {
		return "", nil, ErrNoCollectionTypes
	}

	seen := make(map[CollectionType]struct{}, len(collectionTypes))
	normalizedTypes := make([]CollectionType, 0, len(collectionTypes))
	for _, collectionType := range collectionTypes {
		if !validCollectionType(collectionType) {
			return "", nil, ErrInvalidCollectionType
		}
		if _, ok := seen[collectionType]; ok {
			continue
		}
		seen[collectionType] = struct{}{}
		normalizedTypes = append(normalizedTypes, collectionType)
	}
	sort.Slice(normalizedTypes, func(i, j int) bool {
		return normalizedTypes[i] < normalizedTypes[j]
	})
	return normalizedID, normalizedTypes, nil
}

type pageJob struct {
	params fetchParams
}

type sourceCoordinate struct {
	collectionType CollectionType
	offset         int
	itemIndex      int
}

type locatedSubject struct {
	subject *Subject
	source  sourceCoordinate
}

type locatedPage struct {
	items []locatedSubject
}

func locatePage(page *PageResult, collectionType CollectionType) locatedPage {
	items := make([]locatedSubject, 0, len(page.Data))
	for itemIndex, subject := range page.Data {
		items = append(items, locatedSubject{
			subject: subject,
			source: sourceCoordinate{
				collectionType: collectionType,
				offset:         page.Offset,
				itemIndex:      itemIndex,
			},
		})
	}
	return locatedPage{items: items}
}

func (c *Client) fetchRemainingPages(ctx context.Context, jobs []pageJob) ([]locatedPage, error) {
	results := make([]locatedPage, len(jobs))
	if len(jobs) == 0 {
		return results, nil
	}

	workerCount := min(len(jobs), c.concurrencyLimit)
	group, groupCtx := errgroup.WithContext(ctx)
	var claimMu sync.Mutex
	nextJob := 0

	for range workerCount {
		group.Go(func() error {
			for {
				if err := groupCtx.Err(); err != nil {
					return nil
				}

				claimMu.Lock()
				if nextJob >= len(jobs) {
					claimMu.Unlock()
					return nil
				}
				index := nextJob
				nextJob++
				claimMu.Unlock()

				page, err := c.doRequest(groupCtx, jobs[index].params)
				if err != nil {
					return err
				}
				results[index] = locatePage(page, jobs[index].params.CollectionType)
			}
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, networkErrorForContext(err)
	}
	return results, nil
}

type subjectKey struct {
	subjectType SubjectType
	subjectID   int
	collection  CollectionType
}

func canonicalSubjects(items []locatedSubject) []*Subject {
	winners := make(map[subjectKey]locatedSubject, len(items))
	for _, item := range items {
		key := subjectKey{
			subjectType: item.subject.SubjectType,
			subjectID:   item.subject.SubjectID,
			collection:  item.subject.Type,
		}
		current, exists := winners[key]
		if !exists || coordinateLess(item.source, current.source) {
			winners[key] = item
		}
	}

	result := make([]*Subject, 0, len(winners))
	for _, item := range winners {
		result = append(result, item.subject)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SubjectType != result[j].SubjectType {
			return result[i].SubjectType < result[j].SubjectType
		}
		if result[i].SubjectID != result[j].SubjectID {
			return result[i].SubjectID < result[j].SubjectID
		}
		return result[i].Type < result[j].Type
	})
	if result == nil {
		result = make([]*Subject, 0)
	}
	return result
}

func coordinateLess(left, right sourceCoordinate) bool {
	if left.collectionType != right.collectionType {
		return left.collectionType < right.collectionType
	}
	if left.offset != right.offset {
		return left.offset < right.offset
	}
	return left.itemIndex < right.itemIndex
}

func clamp(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func maxInt() int {
	return int(^uint(0) >> 1)
}
