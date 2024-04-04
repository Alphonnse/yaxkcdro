package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Alphonnse/yaxkcdro/pkg/database/convertor"
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)

type DatabaseClient struct {
	pathToDBFile string
	comicsInfo   map[int]dbModel.DBComicsInfo
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

func (db *DatabaseClient) FindLastDownloadedComic() (int, error) {
	var lastDownloadedComic int
	for i := range db.comicsInfo {
		if i > lastDownloadedComic {
			lastDownloadedComic = i
		}
	}

	return lastDownloadedComic, nil
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
	db.comicsInfo[comicsInfo.Num] = *convertor.FromGlobalToDBComicsInfo(comicsInfo)

	marshledComics, err := json.MarshalIndent(db.comicsInfo, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON data: %s", err.Error())
	}

	err = os.WriteFile(db.pathToDBFile, marshledComics, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON file: %s", err.Error())
	}

	return nil
}
