package words

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Alphonnse/yaxkcdro/pkg/models"
	"github.com/Alphonnse/yaxkcdro/pkg/words/convertor"
	stemmerModel "github.com/Alphonnse/yaxkcdro/pkg/words/models"
	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type Stemmer struct {
}

func NewStemmer(path string) *Stemmer {
	stopwords.LoadStopWordsFromFile(path, "en", "\n")
	return &Stemmer{}
}

func (*Stemmer) StemQueryText(text string) ([]models.StemmedWord, error) {
	stemmedSentence, err := stemSentence(text)
	if err != nil {
		return nil, fmt.Errorf("failed to stem sentence: %s", err.Error())
	}

	return convertor.FromStemmerToGlobalKeywords(stemmedSentence), nil
}

func (*Stemmer) StemComicsDesc(transcript, alt string) ([]models.StemmedWord, error) {
	pattern := `\{\{.*?\}\}`
	re := regexp.MustCompile(pattern)
	transcript = re.ReplaceAllString(transcript, "")

	wholeSentence := fmt.Sprintf("%s %s", alt, transcript)

	stemmedSentence, err := stemSentence(wholeSentence)
	if err != nil {
		return nil, fmt.Errorf("failed to stem sentence: %s", err.Error())
	}

	return convertor.FromStemmerToGlobalKeywords(stemmedSentence), nil
}

func stemSentence(str string) ([]stemmerModel.StemmedWord, error) {
	strWithoutStopwords := stopwords.CleanString(str, "en", false)

	wordsToStem := strings.Fields(strWithoutStopwords)
	stemmedWords := make([]stemmerModel.StemmedWord, 0, len(wordsToStem))

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
			stemmedWords = append(stemmedWords, stemmerModel.StemmedWord{
				Word:  stemmedWord,
				Count: 1,
			})
		}
	}

	return stemmedWords, nil
}
