package repository

import "github.com/Alphonnse/yaxkcdro/internal/models"

type RepositoryService interface {
	GetSavedComics() (map[int]bool, error)
	SaveComics(models.ComicsInfoUC) error
}
