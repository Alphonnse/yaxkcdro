package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

func main() {
	str, err := readArgs()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	loadStopWords("stopwords.txt")
	stemmedSentence, err := stemSentence(str)
	if err != nil {
		log.Fatal(err)
	}

	for _, word := range stemmedSentence {
		fmt.Printf("%s ", word)
	}
}

func readArgs() (string, error) {
	var str string
	flag.StringVar(&str, "s", "", "Sentence to be stemmed")
	flag.Parse()

	if str == "" {
		return "", fmt.Errorf("Error parsing flags")
	}

	return str, nil
}

func loadStopWords(path string) {
	stopwords.LoadStopWordsFromFile(path, "en", "\n")
}

func stemSentence(str string) ([]string, error) {
	strWithoutStopwords := stopwords.CleanString(str, "en", false)

	wordFreq := make(map[string]bool)
	result := make([]string, len(wordFreq))

	words := strings.Fields(strWithoutStopwords)
	for _, word := range words {
		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			return nil, fmt.Errorf("Failed to stem word %s: %v", word, err)
		}
		if !wordFreq[stemmedWord] {
			wordFreq[stemmedWord] = true
			result = append(result, stemmedWord)
		}
	}

	return result, nil
}
