package stemmer

import (
	"fmt"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type Stemmer struct {
}

func NewStemmer(path string) *Stemmer {
	stopwords.LoadStopWordsFromFile(path, "en", "\n")
	return &Stemmer{}
}

func (*Stemmer) StemSentence(str string) ([]StemmedWord, error) {
	strWithoutStopwords := stopwords.CleanString(str, "en", false)

	wordsToStem := strings.Fields(strWithoutStopwords)
	stemmedWords := make([]StemmedWord, 0, len(wordsToStem))

	for _, word := range wordsToStem {
		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			return nil, fmt.Errorf("failed to stem word %s: %s", word, err.Error())
		}

		found := false
		for i, stemmedBefore := range stemmedWords {
			if stemmedBefore.Word == stemmedWord {
				stemmedWords[i].Count++
				found = true
				break
			}
		}
		if !found {
			stemmedWords = append(stemmedWords, StemmedWord{
				Word:  stemmedWord,
				Count: 1,
			})
		}
	}

	return stemmedWords, nil
}
