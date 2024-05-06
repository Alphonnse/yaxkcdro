package xkcd

import "github.com/Alphonnse/yaxkcdro/internal/models"

type XkcdService interface {
	GetComicsCountOnResource() (int, error)
	GetComicsFromResource(comicsNumber int) (*models.ComicsInfoUC, error)
}
