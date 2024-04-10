package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Alphonnse/yaxkcdro/pkg/database/convertor"
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)

type DatabaseClient struct {
	pathToDBFile          string
	comicsInfo            map[int]dbModel.DBComicsInfo
	chunkSize             int
	counterInsertionCalls int
	comicsCount           int
	mu                    sync.RWMutex
}

func NewDatabaseClient(pathToDBFile string) (*DatabaseClient, error) {
	client := &DatabaseClient{
		pathToDBFile: pathToDBFile,
		comicsInfo:   make(map[int]dbModel.DBComicsInfo),
	}

	data, err := os.ReadFile(pathToDBFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading JSON file: %s", err.Error())
	}

	err = json.Unmarshal(data, &client.comicsInfo)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON data: %v", err.Error())
	}

	return client, nil
}
func (db *DatabaseClient) SetChunkSize(chunkSize int, comicsCount int) {
	db.chunkSize = chunkSize
	db.comicsCount = comicsCount
}

func (db *DatabaseClient) GetInstalledComics() map[int]bool {
	mapa := make(map[int]bool, len(db.comicsInfo))
	for i := range db.comicsInfo {
		mapa[i] = true
	}
	return mapa
}

func (db *DatabaseClient) GetComicsInfo(from, count int) []globalModel.ComicInfoToOtput {
	comics := make([]globalModel.ComicInfoToOtput, 0, count)

	for i := from; i < from+count; i++ {
		if comicInfo, err := db.comicsInfo[i]; err == true {
			comics = append(comics, globalModel.ComicInfoToOtput{
				Num:      i,
				Img:      comicInfo.URL,
				Keywords: comicInfo.Keywords,
			})
		}
	}

	return comics
}

func (db *DatabaseClient) InsertComicsIntoDB(comicsInfo globalModel.ComicInfoGlobal) error {
	db.counterInsertionCalls++
	db.mu.Lock()
	defer db.mu.Unlock()

	db.comicsInfo[comicsInfo.Num] = *convertor.FromGlobalToDBComicsInfo(comicsInfo)

	if db.counterInsertionCalls == db.chunkSize || db.comicsCount-len(db.comicsInfo) == 0 {
		marshledComics, err := json.MarshalIndent(db.comicsInfo, "", "	")
		if err != nil {
			return fmt.Errorf("Error marshaling JSON data: %s", err.Error())
		}

		err = os.WriteFile(db.pathToDBFile, marshledComics, 0644)
		if err != nil {
			return fmt.Errorf("Error writing JSON file: %s", err.Error())
		}
		db.counterInsertionCalls = 0
	}

	return nil
}
