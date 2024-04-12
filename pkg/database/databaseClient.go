package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Alphonnse/yaxkcdro/pkg/database/convertor"
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)

type DatabaseClient struct {
	// indexes
	pathToDBFile    string
	pathToIndexFile string
	// comics
	comicsInfo map[int]dbModel.DBComicsInfo
	indexTable map[string][]dbModel.IndexModel
	// db manging info
	chunkSize             int
	counterInsertionCalls int
	comicsCount           int
	mu                    sync.RWMutex
}

func NewDatabaseClient(pathToDBFile, pathToIndexFile string) (*DatabaseClient, error) {
	client := &DatabaseClient{
		pathToDBFile:    pathToDBFile,
		comicsInfo:      make(map[int]dbModel.DBComicsInfo),
		pathToIndexFile: pathToIndexFile,
		indexTable:      make(map[string][]dbModel.IndexModel),
	}

	err := client.readFile(false, pathToDBFile)
	if err != nil {
		return nil, fmt.Errorf("databse file error: %s", err.Error())
	}

	err = client.readFile(true, pathToIndexFile)
	if err != nil {
		return nil, fmt.Errorf("index file error: %s", err.Error())
	}

	return client, nil
}

func (db *DatabaseClient) readFile(index bool, pathToFile string) error {
	_, err := os.Stat(pathToFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File not found, creating new one")
			file, err := os.Create(pathToFile)
			if err != nil {
				return fmt.Errorf("can not create JSON file: %s", err.Error())
			}
			_, err = file.Write([]byte("{}"))
			if err != nil {
				return fmt.Errorf("can not create JSON file: %s", err.Error())
			}
		} else {
			return fmt.Errorf("can not read JSON file: %s", err.Error())
		}
	}

	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return fmt.Errorf("can not read JSON file: %s", err.Error())
	}

	if index {
		err = json.Unmarshal(data, &db.indexTable)
		if err != nil {
			return fmt.Errorf("can not unmarshall JSON data (blank json should contain {}): %v", err.Error())
		}
		return nil
	}

	err = json.Unmarshal(data, &db.comicsInfo)
	if err != nil {
		return fmt.Errorf("can not unmarshall JSON data (blank json should contain {}): %v", err.Error())
	}
	return nil
}

func (db *DatabaseClient) SetChunkSize(chunkSize int, comicsCount int) {
	db.chunkSize = chunkSize
	db.comicsCount = comicsCount
}

func (db *DatabaseClient) InsertComicsIntoFiles(comicsInfo globalModel.ComicInfoGlobal) error {
	db.counterInsertionCalls++
	db.mu.Lock()
	defer db.mu.Unlock()

	db.comicsInfo[comicsInfo.Num] = *convertor.FromGlobalToDBComicsInfo(comicsInfo)
	db.processIndex(comicsInfo.Num, convertor.FromGlobalToDBComicsInfo(comicsInfo))

	if db.counterInsertionCalls == db.chunkSize || db.comicsCount-len(db.comicsInfo) == 0 {
		// write comics
		marshledComics, err := json.MarshalIndent(db.comicsInfo, "", "	")
		if err != nil {
			return fmt.Errorf("can not marshall JSON data: %s", err.Error())
		}

		err = os.WriteFile(db.pathToDBFile, marshledComics, 0644)
		if err != nil {
			return fmt.Errorf("can not write JSON file: %s", err.Error())
		}

		// write index
		marshledIndex, err := json.MarshalIndent(db.indexTable, "", "	")
		if err != nil {
			return fmt.Errorf("can not marshall JSON data: %s", err.Error())
		}

		err = os.WriteFile(db.pathToIndexFile, marshledIndex, 0644)
		if err != nil {
			return fmt.Errorf("can not write JSON file: %s", err.Error())
		}
		db.counterInsertionCalls = 0
	}

	return nil
}

func (db *DatabaseClient) processIndex(comicsID int, comics *dbModel.DBComicsInfo) {
	for _, comicsWord := range comics.Keywords {
		db.indexTable[comicsWord.Word] = append(db.indexTable[comicsWord.Word], dbModel.IndexModel{
			ComicsID: comicsID,
			Weight:   comicsWord.Count,
		})	
	}
}

// with the regular stemmer that doesnt remove dubls
// func (db *DatabaseClient) processIndex(comicsID int, comics *dbModel.DBComicsInfo) {
// 	for _, word := range comics.Keywords {
// 		if _, ok := db.indexTable[word]; ok {
// 			found := false
// 			for i, indexModel := range db.indexTable[word] {
// 				if indexModel.ComicsID == comicsID {
// 					db.indexTable[word][i].Weight++
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				db.indexTable[word] = append(db.indexTable[word], dbModel.IndexModel{
// 					ComicsID: comicsID,
// 					Weight:   1,
// 				})
// 			}
// 		} else {
// 			db.indexTable[word] = append(db.indexTable[word], dbModel.IndexModel{
// 				ComicsID: comicsID,
// 				Weight:   1,
// 			})
// 		}
// 	}
// }

func (db *DatabaseClient) GetWhatComicsAreInstalled() map[int]bool {
	mapa := make(map[int]bool, len(db.comicsInfo))
	for i := range db.comicsInfo {
		mapa[i] = true
	}
	return mapa
}

// func (db *DatabaseClient) FindComicsByStringUsingIndex(word string) map[int]dbModel.DBComicsInfo {
// }
