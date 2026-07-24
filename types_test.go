package collection

import (
	"errors"
	"testing"
	"time"
)

func TestCompleteCollectionDTOAndMutationIsolation(t *testing.T) {
	item := fixtureCollection(42, SubjectTypeAnime, CollectionTypeDone)
	body := fixturePage(1, 50, 0, item)
	params := fetchParams{
		UserID:         "uid",
		SubjectType:    SubjectTypeAnime,
		CollectionType: CollectionTypeDone,
		Limit:          50,
		Offset:         0,
	}

	first, err := decodePage(body, params)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Data) != 1 {
		t.Fatalf("data length = %d", len(first.Data))
	}
	subject := first.Data[0]
	expectedTime, _ := time.Parse(time.RFC3339, "2026-07-24T12:34:56Z")
	if subject.ID != 42 || subject.SubjectID != 42 ||
		subject.SubjectType != SubjectTypeAnime ||
		subject.Type != CollectionTypeDone ||
		subject.Name != "subject-42" ||
		subject.NameCn != "条目-42" ||
		subject.Rate != 8 ||
		subject.Comment != "fixture comment" ||
		!subject.UpdatedAt.Equal(expectedTime) ||
		subject.VolStatus != 2 ||
		subject.EpStatus != 7 ||
		subject.Private ||
		len(subject.Tags) != 2 {
		t.Fatalf("decoded subject = %#v", subject)
	}

	first.Data[0].Tags[0] = "mutated"
	first.Data[0].Name = "mutated"
	second, err := decodePage(body, params)
	if err != nil {
		t.Fatal(err)
	}
	if second.Data[0].Tags[0] != "one" || second.Data[0].Name != "subject-42" {
		t.Fatalf("second result shared caller mutation: %#v", second.Data[0])
	}
}

func TestRequiredTagsAndOptionalCommentSubject(t *testing.T) {
	params := fetchParams{
		SubjectType:    SubjectTypeAnime,
		CollectionType: CollectionTypeDone,
		Limit:          50,
	}
	tests := []struct {
		name        string
		mutate      func(map[string]any)
		wantComment string
		wantName    string
		wantNameCn  string
	}{
		{
			name: "empty required tags",
			mutate: func(item map[string]any) {
				item["tags"] = []string{}
			},
			wantComment: "fixture comment",
			wantName:    "subject-1",
			wantNameCn:  "条目-1",
		},
		{
			name: "omitted comment",
			mutate: func(item map[string]any) {
				delete(item, "comment")
			},
			wantName:   "subject-1",
			wantNameCn: "条目-1",
		},
		{
			name: "omitted subject",
			mutate: func(item map[string]any) {
				delete(item, "subject")
			},
			wantComment: "fixture comment",
		},
		{
			name: "both optional fields omitted",
			mutate: func(item map[string]any) {
				delete(item, "comment")
				delete(item, "subject")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := fixtureCollection(1, SubjectTypeAnime, CollectionTypeDone)
			test.mutate(item)
			page, err := decodePage(fixturePage(1, 50, 0, item), params)
			if err != nil {
				t.Fatal(err)
			}
			subject := page.Data[0]
			if subject.ID != 1 || subject.SubjectID != 1 ||
				subject.Comment != test.wantComment ||
				subject.Name != test.wantName ||
				subject.NameCn != test.wantNameCn {
				t.Fatalf("subject = %#v", subject)
			}
			if subject.Tags == nil {
				t.Fatal("tags is nil")
			}
		})
	}
}

func TestCollectionProtocolValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "missing updated at", mutate: func(item map[string]any) { delete(item, "updated_at") }},
		{name: "missing required tags", mutate: func(item map[string]any) { delete(item, "tags") }},
		{name: "null required tags", mutate: func(item map[string]any) { item["tags"] = nil }},
		{name: "malformed required tags", mutate: func(item map[string]any) { item["tags"] = "one" }},
		{name: "null required tag item", mutate: func(item map[string]any) {
			item["tags"] = []any{"one", nil}
		}},
		{name: "null optional comment", mutate: func(item map[string]any) { item["comment"] = nil }},
		{name: "malformed optional comment", mutate: func(item map[string]any) { item["comment"] = 1 }},
		{name: "null optional subject", mutate: func(item map[string]any) { item["subject"] = nil }},
		{name: "malformed optional subject", mutate: func(item map[string]any) { item["subject"] = "subject" }},
		{name: "incomplete optional subject", mutate: func(item map[string]any) {
			delete(item["subject"].(map[string]any), "name_cn")
		}},
		{name: "zero id", mutate: func(item map[string]any) {
			item["subject_id"] = 0
			item["subject"].(map[string]any)["id"] = 0
		}},
		{name: "id mismatch", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["id"] = 99
		}},
		{name: "subject type mismatch", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = int(SubjectTypeBook)
		}},
		{name: "invalid subject type", mutate: func(item map[string]any) { item["subject_type"] = 5 }},
		{name: "wrong requested subject type", mutate: func(item map[string]any) {
			item["subject_type"] = int(SubjectTypeBook)
			item["subject"].(map[string]any)["type"] = int(SubjectTypeBook)
		}},
		{name: "invalid collection type", mutate: func(item map[string]any) { item["type"] = 6 }},
		{name: "wrong requested collection type", mutate: func(item map[string]any) {
			item["type"] = int(CollectionTypeWish)
		}},
		{name: "negative rate", mutate: func(item map[string]any) { item["rate"] = -1 }},
		{name: "high rate", mutate: func(item map[string]any) { item["rate"] = 11 }},
		{name: "negative volumes", mutate: func(item map[string]any) { item["vol_status"] = -1 }},
		{name: "negative episodes", mutate: func(item map[string]any) { item["ep_status"] = -1 }},
		{name: "invalid timestamp", mutate: func(item map[string]any) { item["updated_at"] = "yesterday" }},
		{name: "missing private", mutate: func(item map[string]any) { delete(item, "private") }},
	}

	params := fetchParams{
		SubjectType:    SubjectTypeAnime,
		CollectionType: CollectionTypeDone,
		Limit:          50,
		Offset:         0,
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := fixtureCollection(1, SubjectTypeAnime, CollectionTypeDone)
			test.mutate(item)
			_, err := decodePage(fixturePage(1, 50, 0, item), params)
			var protocolErr *ProtocolError
			if !errors.As(err, &protocolErr) || !errors.Is(err, ErrProtocol) {
				t.Fatalf("error = %T %v", err, err)
			}
		})
	}
}

func TestPageMetadataAndStrictJSONValidation(t *testing.T) {
	item := fixtureCollection(1, SubjectTypeAnime, CollectionTypeDone)
	params := fetchParams{
		SubjectType:    SubjectTypeAnime,
		CollectionType: CollectionTypeDone,
		Limit:          1,
		Offset:         3,
	}
	protocolBodies := [][]byte{
		fixturePage(-1, 1, 3),
		fixturePage(1_000_001, 1, 3),
		fixturePage(0, 2, 3),
		fixturePage(0, 1, 2),
		fixturePage(2, 1, 3, item, item),
		[]byte(`{"data":[],"limit":1,"offset":3}`),
		[]byte(`{"total":0,"limit":1,"offset":3}`),
	}
	for index, body := range protocolBodies {
		_, err := decodePage(body, params)
		if !errors.Is(err, ErrProtocol) {
			t.Fatalf("protocol case %d: %T %v", index, err, err)
		}
	}

	decodeBodies := [][]byte{
		[]byte(`not json`),
		append(fixturePage(0, 1, 3), []byte(` {}`)...),
	}
	for index, body := range decodeBodies {
		_, err := decodePage(body, params)
		var decodeErr *DecodeError
		if !errors.As(err, &decodeErr) || !errors.Is(err, ErrDecode) {
			t.Fatalf("decode case %d: %T %v", index, err, err)
		}
	}
}
