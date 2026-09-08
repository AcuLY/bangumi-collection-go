package collection

import (
	"bytes"
	"encoding/json"
	"time"
)

// SubjectType is a Bangumi subject category.
type SubjectType int

const (
	SubjectTypeBook  SubjectType = 1
	SubjectTypeAnime SubjectType = 2
	SubjectTypeMusic SubjectType = 3
	SubjectTypeGame  SubjectType = 4
	SubjectTypeReal  SubjectType = 6
)

// CollectionType is a user's collection state.
type CollectionType int

const (
	CollectionTypeWish    CollectionType = 1
	CollectionTypeDone    CollectionType = 2
	CollectionTypeDoing   CollectionType = 3
	CollectionTypeOnHold  CollectionType = 4
	CollectionTypeDropped CollectionType = 5
)

// Subject is one complete public collection record.
//
// ID is retained as a compatibility alias and always equals SubjectID.
type Subject struct {
	ID          int            `json:"id"`
	SubjectID   int            `json:"subject_id"`
	SubjectType SubjectType    `json:"subject_type"`
	Type        CollectionType `json:"type"`
	Name        string         `json:"name"`
	NameCn      string         `json:"name_cn"`
	Rate        int            `json:"rate"`
	Comment     string         `json:"comment"`
	Tags        []string       `json:"tags"`
	UpdatedAt   time.Time      `json:"updated_at"`
	VolStatus   int            `json:"vol_status"`
	EpStatus    int            `json:"ep_status"`
	Private     bool           `json:"private"`
}

// PageResult is one validated upstream page. Data preserves upstream order.
type PageResult struct {
	Data   []*Subject
	Total  int
	Limit  int
	Offset int
}

type fetchParams struct {
	UserID         string
	SubjectType    SubjectType
	CollectionType CollectionType
	Offset         int
	Limit          int
}

type wirePage struct {
	Data   *[]wireCollection `json:"data"`
	Total  *int64            `json:"total"`
	Limit  *int64            `json:"limit"`
	Offset *int64            `json:"offset"`
}

type wireCollection struct {
	UpdatedAt   *string         `json:"updated_at"`
	Comment     json.RawMessage `json:"comment"`
	Tags        json.RawMessage `json:"tags"`
	Subject     json.RawMessage `json:"subject"`
	SubjectID   *int64          `json:"subject_id"`
	VolStatus   *int64          `json:"vol_status"`
	EpStatus    *int64          `json:"ep_status"`
	SubjectType *int64          `json:"subject_type"`
	Type        *int64          `json:"type"`
	Rate        *int64          `json:"rate"`
	Private     *bool           `json:"private"`
}

type wireSubject struct {
	ID     *int64  `json:"id"`
	Type   *int64  `json:"type"`
	Name   *string `json:"name"`
	NameCn *string `json:"name_cn"`
}

func validSubjectType(value SubjectType) bool {
	switch value {
	case SubjectTypeBook, SubjectTypeAnime, SubjectTypeGame, SubjectTypeMusic, SubjectTypeReal:
		return true
	default:
		return false
	}
}

func validCollectionType(value CollectionType) bool {
	return value >= CollectionTypeWish && value <= CollectionTypeDropped
}

func (item wireCollection) toSubject(expectedSubject SubjectType, expectedCollection CollectionType) (*Subject, error) {
	if item.UpdatedAt == nil || item.SubjectID == nil ||
		item.VolStatus == nil || item.EpStatus == nil ||
		item.SubjectType == nil || item.Type == nil || item.Rate == nil ||
		item.Private == nil {
		return nil, newProtocolError()
	}

	subjectID := *item.SubjectID
	subjectType := SubjectType(*item.SubjectType)
	collectionType := CollectionType(*item.Type)
	if subjectID <= 0 ||
		!validSubjectType(subjectType) ||
		subjectType != expectedSubject ||
		!validCollectionType(collectionType) ||
		collectionType != expectedCollection ||
		*item.Rate < 0 || *item.Rate > 10 ||
		*item.VolStatus < 0 || *item.EpStatus < 0 ||
		!fitsInt(subjectID) || !fitsInt(*item.VolStatus) || !fitsInt(*item.EpStatus) {
		return nil, newProtocolError()
	}

	comment, err := decodeOptionalString(item.Comment)
	if err != nil {
		return nil, err
	}
	tags, err := decodeRequiredTags(item.Tags)
	if err != nil {
		return nil, err
	}
	name, nameCn, err := decodeOptionalSubject(item.Subject, subjectID)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, *item.UpdatedAt)
	if err != nil {
		return nil, newProtocolError()
	}

	return &Subject{
		ID:          int(subjectID),
		SubjectID:   int(subjectID),
		SubjectType: subjectType,
		Type:        collectionType,
		Name:        name,
		NameCn:      nameCn,
		Rate:        int(*item.Rate),
		Comment:     comment,
		Tags:        tags,
		UpdatedAt:   updatedAt,
		VolStatus:   int(*item.VolStatus),
		EpStatus:    int(*item.EpStatus),
		Private:     *item.Private,
	}, nil
}

func decodeOptionalString(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return "", nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", newProtocolError()
	}
	return value, nil
}

func decodeRequiredTags(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, newProtocolError()
	}
	var encodedTags []json.RawMessage
	if err := json.Unmarshal(raw, &encodedTags); err != nil || encodedTags == nil {
		return nil, newProtocolError()
	}
	tags := make([]string, len(encodedTags))
	for index, encoded := range encodedTags {
		if bytes.Equal(bytes.TrimSpace(encoded), []byte("null")) ||
			json.Unmarshal(encoded, &tags[index]) != nil {
			return nil, newProtocolError()
		}
	}
	return tags, nil
}

func decodeOptionalSubject(
	raw json.RawMessage,
	subjectID int64,
) (string, string, error) {
	if len(raw) == 0 {
		return "", "", nil
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return "", "", newProtocolError()
	}

	var subject wireSubject
	if err := json.Unmarshal(raw, &subject); err != nil ||
		subject.ID == nil || subject.Type == nil ||
		subject.Name == nil || subject.NameCn == nil ||
		*subject.ID != subjectID ||
		!validSubjectType(SubjectType(*subject.Type)) {
		return "", "", newProtocolError()
	}
	// Subject metadata may have been reclassified independently of the
	// collection record. The public SubjectType keeps the top-level value.
	return *subject.Name, *subject.NameCn, nil
}

func fitsInt(value int64) bool {
	converted := int(value)
	return int64(converted) == value
}
