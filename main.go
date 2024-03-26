package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

func main() {
	str := readSFromArgs()

	for _, word := range stemSentence(str) {
		fmt.Printf("%s ", word)
	}
}

func readSFromArgs() string {
	var str string
	cliArgs := os.Args

	if len(cliArgs) == 3 {
		if cliArgs[1] == "-s" {
			str = cliArgs[2]
		} else {
			log.Fatal("Wrong key. Please use -s key to specify a sentence")
		}
	} else {
		log.Fatal("Please use -s key only to specify a sentence")
	}
	return str
}

func stemSentence(str string) []string {
	stopwords.LoadStopWordsFromFile("stopwords.txt", "en", "\n")
	strWithoutStopwords := stopwords.CleanString(str, "en", true)

	stemmed, _ := snowball.Stem(strWithoutStopwords, "english", true)

	return removeDuplicates(stemmed)
}

func removeDuplicates(str string) []string {
	words := make(map[string]bool)
	// i prefer the slice of strings instead of string
	// because concatenation of strings is too expensive in memory
	result := make([]string, len(words))

	for _, word := range strings.Split(str, " ") {
		if !words[word] {
			words[word] = true
			result = append(result, word)
		}
	}

	return result
}
