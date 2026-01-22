package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AcuLY/bangumi-collection-go"
)

func main() {
	client := collection.NewClient("AcuL/bangumi-collection-go")

	subjects, err := client.Fetch(
		context.Background(),
		"lucay126",                          
		collection.SubjectTypeAnime,    
		collection.CollectionTypeDoing,
		collection.CollectionTypeDone, 
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("共 %d 部动画\n\n", len(subjects))
	for _, s := range subjects {
		name := s.NameCn
		if name == "" {
			name = s.Name
		}
		fmt.Printf("ID: %d | %s | 评分: %d | 标签: %v\n", s.ID, name, s.Rate, s.Tags)
	}
}
