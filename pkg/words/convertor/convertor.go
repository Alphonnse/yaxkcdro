package convertor

import (
	"github.com/Alphonnse/yaxkcdro/pkg/models"
	stemmerModel "github.com/Alphonnse/yaxkcdro/pkg/words/models"
)

func FromStemmerToGlobalKeywords(stemmedWords []stemmerModel.StemmedWord) []models.StemmedWord {
	var keywords []models.StemmedWord
	for _, word := range stemmedWords {
		keywords = append(keywords, models.StemmedWord{
			Word:  word.Word,
			Count: word.Count,
		})
	}
	return keywords
}
