package words

import "github.com/Alphonnse/yaxkcdro/pkg/models"

type StemmerService interface {
	StemQueryText(text string) ([]models.StemmedWord, error)
	StemComicsDesc(transcript, alt string) ([]models.StemmedWord, error)
}
