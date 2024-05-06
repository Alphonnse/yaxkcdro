package stemmeradapter

import (
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/stemmer"
)

func FromStemmerStemSentenceToGlobal(stemmedWords []stemmer.StemmedWord) []models.StemmedWordUC {
	out := make([]models.StemmedWordUC, 0, len(stemmedWords))

	for _, keyword := range stemmedWords {
		out = append(out, models.StemmedWordUC{
			Word:  keyword.Word,
			Count: keyword.Count,
		})
	}

	return out
}
