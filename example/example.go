package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	collection "github.com/AcuLY/bangumi-collection-go"
)

func main() {
	// The official mapping is SubjectTypeMusic=3 and SubjectTypeGame=4.
	// These names correct the reversed values in the untagged prototype.
	client := collection.NewClient(
		"AcuL/bangumi-collection-go-example",
		collection.WithConcurrencyLimit(10),
		collection.WithRateLimit(3, 1),
		collection.WithRequestTimeout(30*time.Second),
		collection.WithMaxRetryDelay(30*time.Second),
	)

	subjects, err := client.Fetch(
		context.Background(),
		"lucay126",
		collection.SubjectTypeAnime,
		collection.CollectionTypeDoing,
		collection.CollectionTypeDone,
	)
	if err != nil {
		var httpErr *collection.HTTPError
		if errors.As(err, &httpErr) {
			log.Printf("Bangumi HTTP status: %d", httpErr.StatusCode)
		}
		if errors.Is(err, collection.ErrRateLimited) {
			log.Fatal("Bangumi request was rate limited")
		}
		log.Fatal(err)
	}

	fmt.Printf("共 %d 部动画\n\n", len(subjects))
	for _, subject := range subjects {
		name := subject.NameCn
		if name == "" {
			name = subject.Name
		}
		fmt.Printf(
			"ID: %d | %s | 收藏状态: %d | 评分: %d | 评论: %q | 标签: %v | 更新时间: %s\n",
			subject.SubjectID,
			name,
			subject.Type,
			subject.Rate,
			subject.Comment,
			subject.Tags,
			subject.UpdatedAt.Format(time.RFC3339),
		)
	}
}
