package stemmer

import (
	"github.com/Alphonnse/yaxkcdro/internal/models"
)

type StemmerService interface {
	StemComicsDescription(title, transcript, alt string) ([]models.StemmedWordUC, error)
}
