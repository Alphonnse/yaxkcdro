package words

import "github.com/Alphonnse/yaxkcdro/pkg/models"

type StemmerService interface {
	// Stem(comicsInfo models.ComicInfoGlobal) (*models.ComicInfoGlobal, error)
	StemComicsDesc(transcript, alt string) ([]models.Keyword, error)
}
