package database

import (
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
	"testing"
)

var queryWords = []globalModel.StemmedWord{
	{Word: "follow", Count: 1},
	{Word: "question", Count: 1},
}

func BenchmarkFindComicsByStringNotUsingIndex(b *testing.B) {
	client, _ := NewDatabaseClient("../../database.json", "../../index.json")

	for i := 0; i < b.N; i++ {
		_ = client.FindComicsByStringNotUsingIndex(queryWords)
		// _ = client.FindComicsByStringUsingIndex(queryWords)
	}
}

func BenchmarkFindComicsByStringUsingIndex(b *testing.B) {
	client, _ := NewDatabaseClient("../../database.json", "../../index.json")

	for i := 0; i < b.N; i++ {
		_ = client.FindComicsByStringUsingIndex(queryWords)
	}
}
