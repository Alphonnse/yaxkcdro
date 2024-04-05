package database

import globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"

type DatabaseService interface {
	FindLastDownloadedComic() (int, error)
	GetComicsInfo(from, count int) []globalModel.ComicInfoToOtput
	InsertComicsIntoDB(comicsInfo globalModel.ComicInfoGlobal) error
}
