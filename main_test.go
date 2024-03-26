package main

import (
	"fmt"
	"testing"
)

// benchmarking
func BenchmarkRemoveDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		removeDuplicates("You can find an example here https:github.com/bbalet/gorelated where stopwords package is used in conjunction with SimHash algorithm in order to find a list of related content for a static website generator:You can find an example here https:github.com/bbalet/gorelated where stopwords package is used in conjunction with SimHash algorithm in order to find a list of related content for a static website generator:")
	}
}

// Unit test of stemSentence
func TestStemSentence(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "The quick brown fox jumps over the lazy dog",
			expected: []string{"quick", "brown", "fox", "jumps", "lazy", "dog"},
		},
		{
			input:    "i'll follow you as long as you are following me",
			expected: []string{"follow", "long"},
		},
		{
			input:    "A timestamp server works by taking a hash of a block of items to be timestamped and widely publishing the hash, such as in a newspaper or Usenet post",
			expected: []string{"timestamp", "server", "works" , "taking", "hash", "block", "items", "timestamped", "widely", "publishing", "newspaper", "usenet", "post"},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := stemSentence(testCase.input)
			if len(actual) != len(testCase.expected) {
				t.Errorf("Test case %d: Expected %v words, but got %v", i, len(testCase.expected), len(actual))
			}
			for j, word := range actual {
				if word != testCase.expected[j] {
					t.Errorf("Test case %d: Expected word at index %v to be %v, but got %v", i, j, testCase.expected[j], word)
				}
			}
		})
	}
}
