package database

import globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"

type DatabaseService interface {
	SetChunkSize(chunkSize int, comicsCount int)
	GetInstalledComics() map[int]bool
	InsertComicsIntoDB(comicsInfo globalModel.ComicInfoGlobal) error
}
