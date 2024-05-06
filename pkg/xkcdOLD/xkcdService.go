package xkcdClient

import "github.com/Alphonnse/yaxkcdro/pkg/models"

type XkcsService interface {
	GetComicsFromResource(lastDownloadedComic int) (*models.ComicInfoGlobal, error)
	GetComicsCountOnResource() (int, error)
}
