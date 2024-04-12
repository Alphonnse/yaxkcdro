package database

import (
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
)

type DatabaseService interface {
	SetChunkSize(chunkSize int, comicsCount int)
	GetWhatComicsAreInstalled() map[int]bool
	readFile(index bool, pathToFile string) error
	InsertComicsIntoFiles(comicsInfo globalModel.ComicInfoGlobal) error
	processIndex(comicsID int, comics *dbModel.DBComicsInfo)
}
