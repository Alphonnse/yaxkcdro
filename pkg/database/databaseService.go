package database

import (
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)

type DatabaseService interface {
	SetChunkSize(chunkSize int, comicsCount int)
	GetWhatComicsAreInstalled() map[int]bool
	readFile(index bool, pathToFile string) error
	InsertComicsIntoFiles(comicsInfo globalModel.ComicInfoGlobal) error
	processIndex(comicsID int, comics *dbModel.DBComicsInfo)
	FindComicsByStringUsingIndex(queryWords []globalModel.StemmedWord) []globalModel.ComicInfoToOtput
	FindComicsByStringNotUsingIndex(queryWords []globalModel.StemmedWord) []globalModel.ComicInfoToOtput
}
