package collection

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
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

func TestSupportedNestedSubjectTypesPreserveCollectionMetadata(t *testing.T) {
	params := fetchParams{
		SubjectType:    SubjectTypeAnime,
		CollectionType: CollectionTypeDone,
		Limit:          50,
	}
	tests := []struct {
		name        string
		subjectType SubjectType
	}{
		{name: "book", subjectType: SubjectTypeBook},
		{name: "anime", subjectType: SubjectTypeAnime},
		{name: "music", subjectType: SubjectTypeMusic},
		{name: "game", subjectType: SubjectTypeGame},
		{name: "real", subjectType: SubjectTypeReal},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := fixtureCollection(42, SubjectTypeAnime, CollectionTypeDone)
			item["subject"].(map[string]any)["type"] = int(test.subjectType)
			page, err := decodePage(fixturePage(1, 50, 0, item), params)
			if err != nil {
				t.Fatal(err)
			}
			if len(page.Data) != 1 {
				t.Fatalf("data length = %d", len(page.Data))
			}
			subject := page.Data[0]
			if subject.ID != 42 || subject.SubjectID != 42 ||
				subject.SubjectType != SubjectTypeAnime ||
				subject.Type != CollectionTypeDone ||
				subject.Name != "subject-42" || subject.NameCn != "条目-42" {
				t.Fatalf("decoded subject = %#v", subject)
			}
		})
	}
}

func TestReclassifiedSubjectPageAndAggregate(t *testing.T) {
	t.Run("page", func(t *testing.T) {
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			item := fixtureCollection(631949, SubjectTypeAnime, collectionType)
			item["subject"].(map[string]any)["type"] = int(SubjectTypeReal)
			writeJSON(t, writer, fixturePage(1, limit, offset, item))
		}))

		page, err := client.FetchPage(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
			50,
			0,
		)
		if err != nil {
			t.Fatal(err)
		}
		if page == nil || len(page.Data) != 1 {
			t.Fatalf("page = %#v", page)
		}
		subject := page.Data[0]
		if subject.ID != 631949 || subject.SubjectID != 631949 ||
			subject.SubjectType != SubjectTypeAnime ||
			subject.Type != CollectionTypeDone ||
			subject.Name != "subject-631949" || subject.NameCn != "条目-631949" {
			t.Fatalf("subject = %#v", subject)
		}
	})

	t.Run("aggregate", func(t *testing.T) {
		var calls atomic.Int64
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			calls.Add(1)
			collectionType, limit, offset := requestPageCoordinates(t, request)
			switch offset {
			case 0:
				items := make([]map[string]any, 50)
				for index := range items {
					items[index] = fixtureCollection(50-index, SubjectTypeAnime, collectionType)
				}
				writeJSON(t, writer, fixturePage(51, limit, offset, items...))
			case 50:
				item := fixtureCollection(631949, SubjectTypeAnime, collectionType)
				item["subject"].(map[string]any)["type"] = int(SubjectTypeReal)
				writeJSON(t, writer, fixturePage(51, limit, offset, item))
			default:
				t.Errorf("unexpected offset %d", offset)
				writeJSON(t, writer, fixturePage(0, limit, offset))
			}
		}))

		subjects, err := client.Fetch(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(subjects) != 51 || calls.Load() != 2 {
			t.Fatalf("subjects length = %d, transport calls = %d", len(subjects), calls.Load())
		}
		for index, subject := range subjects {
			wantID := index + 1
			if index == 50 {
				wantID = 631949
			}
			if subject.ID != wantID || subject.SubjectID != wantID ||
				subject.SubjectType != SubjectTypeAnime ||
				subject.Type != CollectionTypeDone ||
				subject.Name != fmt.Sprintf("subject-%d", wantID) ||
				subject.NameCn != fmt.Sprintf("条目-%d", wantID) {
				t.Fatalf("subject at index %d = %#v", index, subject)
			}
		}
	})
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
			name: "null comment",
			mutate: func(item map[string]any) {
				item["comment"] = nil
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

func TestNullableCommentPageAndAggregate(t *testing.T) {
	t.Run("page", func(t *testing.T) {
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			item := fixtureCollection(1, SubjectTypeAnime, collectionType)
			item["comment"] = nil
			writeJSON(t, writer, fixturePage(1, limit, offset, item))
		}))

		page, err := client.FetchPage(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
			50,
			0,
		)
		if err != nil {
			t.Fatal(err)
		}
		if page == nil || len(page.Data) != 1 || page.Data[0].Comment != "" {
			t.Fatalf("page = %#v", page)
		}
	})

	t.Run("aggregate", func(t *testing.T) {
		client, _ := newLoopbackServerClient(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			collectionType, limit, offset := requestPageCoordinates(t, request)
			switch offset {
			case 0:
				withString := fixtureCollection(1, SubjectTypeAnime, collectionType)
				omitted := fixtureCollection(2, SubjectTypeAnime, collectionType)
				delete(omitted, "comment")
				writeJSON(t, writer, fixturePage(51, limit, offset, withString, omitted))
			case 50:
				nullComment := fixtureCollection(3, SubjectTypeAnime, collectionType)
				nullComment["comment"] = nil
				writeJSON(t, writer, fixturePage(51, limit, offset, nullComment))
			default:
				t.Errorf("unexpected offset %d", offset)
				writeJSON(t, writer, fixturePage(0, limit, offset))
			}
		}))

		subjects, err := client.Fetch(
			context.Background(),
			"uid",
			SubjectTypeAnime,
			CollectionTypeDone,
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(subjects) != 3 ||
			subjects[0].Comment != "fixture comment" ||
			subjects[1].Comment != "" ||
			subjects[2].Comment != "" {
			t.Fatalf("subjects = %#v", subjects)
		}
	})
}

func TestNonStringOptionalCommentIsProtocolFailure(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "number", value: 1},
		{name: "object", value: map[string]any{"fixture": "value"}},
		{name: "array", value: []any{"fixture"}},
		{name: "boolean", value: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int64
			client, _ := newLoopbackServerClient(
				t,
				http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					calls.Add(1)
					collectionType, limit, offset := requestPageCoordinates(t, request)
					valid := fixtureCollection(1, SubjectTypeAnime, collectionType)
					invalid := fixtureCollection(2, SubjectTypeAnime, collectionType)
					invalid["comment"] = test.value
					writeJSON(t, writer, fixturePage(2, limit, offset, valid, invalid))
				}),
				WithMaxRetries(3),
			)

			page, err := client.FetchPage(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
				50,
				0,
			)
			var protocolErr *ProtocolError
			if page != nil || !errors.As(err, &protocolErr) || !errors.Is(err, ErrProtocol) {
				t.Fatalf("page=%#v error=%T %v", page, err, err)
			}

			subjects, err := client.Fetch(
				context.Background(),
				"uid",
				SubjectTypeAnime,
				CollectionTypeDone,
			)
			protocolErr = nil
			if subjects != nil || !errors.As(err, &protocolErr) || !errors.Is(err, ErrProtocol) {
				t.Fatalf("subjects=%#v error=%T %v", subjects, err, err)
			}
			if calls.Load() != 2 {
				t.Fatalf("transport calls = %d", calls.Load())
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
		{name: "missing nested subject type", mutate: func(item map[string]any) {
			delete(item["subject"].(map[string]any), "type")
		}},
		{name: "null nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = nil
		}},
		{name: "string nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = "6"
		}},
		{name: "fractional nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = 6.5
		}},
		{name: "boolean nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = true
		}},
		{name: "object nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = map[string]any{"type": 6}
		}},
		{name: "array nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = []int{6}
		}},
		{name: "negative nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = -1
		}},
		{name: "zero nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = 0
		}},
		{name: "unsupported nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = 5
		}},
		{name: "high nested subject type", mutate: func(item map[string]any) {
			item["subject"].(map[string]any)["type"] = 7
		}},
		{name: "invalid subject type", mutate: func(item map[string]any) { item["subject_type"] = 5 }},
		{name: "wrong requested subject type", mutate: func(item map[string]any) {
			item["subject_type"] = int(SubjectTypeBook)
			item["subject"].(map[string]any)["type"] = int(SubjectTypeBook)
		}},
		{name: "wrong requested outer type with matching nested type", mutate: func(item map[string]any) {
			item["subject_type"] = int(SubjectTypeBook)
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
