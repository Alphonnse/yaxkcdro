package words

import "github.com/Alphonnse/yaxkcdro/pkg/models"

type StemmerService interface {
	Stem(comicsInfo models.ComicInfoGlobal) (*models.ComicInfoGlobal, error)
}
