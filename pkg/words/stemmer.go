package words

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Alphonnse/yaxkcdro/pkg/models"
	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type Stemmer struct {
}

func NewStemmer(path string) *Stemmer {
	stopwords.LoadStopWordsFromFile(path, "en", "\n")
	return &Stemmer{}
}

func (*Stemmer) Stem(comicsInfo models.ComicInfoGlobal) (*models.ComicInfoGlobal, error) {
	pattern := `\{\{.*?\}\}`
	re := regexp.MustCompile(pattern)
	comicsInfo.Transcript = re.ReplaceAllString(comicsInfo.Transcript, "")

	wholeSentence := fmt.Sprintf("%s %s", comicsInfo.Alt, comicsInfo.Transcript)

	stemmedSentence, err := stemSentence(wholeSentence)
	if err != nil {
		return nil, fmt.Errorf("failed to stem sentence: %s", err.Error())
	}

	comicsInfo.Keywords = stemmedSentence
	return &comicsInfo, nil
}

func stemSentence(str string) ([]string, error) {
	strWithoutStopwords := stopwords.CleanString(str, "en", false)

	wordFreq := make(map[string]bool)
	result := make([]string, len(wordFreq))

	words := strings.Fields(strWithoutStopwords)
	for _, word := range words {
		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			return nil, fmt.Errorf("failed to stem word %s: %s", word, err.Error())
		}
		if !wordFreq[stemmedWord] {
			wordFreq[stemmedWord] = true
			result = append(result, stemmedWord)
		}
	}

	return result, nil
}
