package stemmer

import (
	"fmt"
	"regexp"

	stemmeradapter "github.com/Alphonnse/yaxkcdro/internal/adapters/stemmerAdapter"
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/stemmer"
)

type StemmerClient struct {
	stemmerClient *stemmer.Stemmer
}

func NewStemmerClient(pathToStropwords string) *StemmerClient {
	return &StemmerClient{
		stemmerClient: stemmer.NewStemmer(pathToStropwords),
	}
}

func (s *StemmerClient) StemComicsDescription(title, transcript, alt string) ([]models.StemmedWordUC, error) {
	pattern := `\{\{.*?\}\}`
	re := regexp.MustCompile(pattern)
	transcript = re.ReplaceAllString(transcript, "")

	wholeSentence := fmt.Sprintf("%s %s %s", title, alt, transcript)

	stemmedSentence, err := s.stemmerClient.StemSentence(wholeSentence)
	if err != nil {
		return nil, fmt.Errorf("failed to stem sentence: %s", err.Error())
	}

	return stemmeradapter.FromStemmerStemSentenceToGlobal(stemmedSentence), nil
}
