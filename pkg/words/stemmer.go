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

// func (*Stemmer) Stem(comicsInfo models.ComicInfoGlobal) (*models.ComicInfoGlobal, error) {
// 	pattern := `\{\{.*?\}\}`
// 	re := regexp.MustCompile(pattern)
// 	comicsInfo.Transcript = re.ReplaceAllString(comicsInfo.Transcript, "")
//
// 	wholeSentence := fmt.Sprintf("%s %s", comicsInfo.Alt, comicsInfo.Transcript)
//
// 	stemmedSentence, err := stemSentence(wholeSentence)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to stem sentence: %s", err.Error())
// 	}
//
// 	comicsInfo.Keywords = stemmedSentence
// 	return &comicsInfo, nil
// }

func (*Stemmer) StemComicsDesc(transcript, alt string) ([]models.Keyword, error) {
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

// func stemSentence(str string) ([]string, error) {
// 	strWithoutStopwords := stopwords.CleanString(str, "en", false)
//
// 	words := strings.Fields(strWithoutStopwords)
// 	result := make([]string, len(words))
//
// 	for _, word := range words {
// 		stemmedWord, err := snowball.Stem(word, "english", true)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to stem word %s: %s", word, err.Error())
// 		}
// 		result = append(result, stemmedWord)
// 	}
//
// 	return result, nil
// }
