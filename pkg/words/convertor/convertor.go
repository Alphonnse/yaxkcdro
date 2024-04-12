package convertor

import (
	"github.com/Alphonnse/yaxkcdro/pkg/models"
	stemmerModel "github.com/Alphonnse/yaxkcdro/pkg/words/models"
)

func FromStemmerToGlobalKeywords(stemmedWords []stemmerModel.StemmedWord) []models.Keyword {
	var keywords []models.Keyword
	for _, word := range stemmedWords {
		keywords = append(keywords, models.Keyword{
			Word:  word.Word,
			Count: word.Count,
		})
	}
	return keywords
}
