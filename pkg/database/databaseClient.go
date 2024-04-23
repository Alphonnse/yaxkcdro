package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
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
		marshledComics, err := json.MarshalIndent(db.comicsInfo, "", "	")
		if err != nil {
			return fmt.Errorf("can not marshall JSON data: %s", err.Error())
		}

		marshledIndex, err := json.MarshalIndent(db.indexTable, "", "	")
		if err != nil {
			return fmt.Errorf("can not marshall JSON data: %s", err.Error())
		}

		err = db.writeFile(db.pathToDBFile, marshledComics)
		if err != nil {
			return fmt.Errorf("can not write JSON DB file: %s", err.Error())
		}

		err = db.writeFile(db.pathToIndexFile, marshledIndex)
		if err != nil {
			return fmt.Errorf("can not write JSON index file: %s", err.Error())
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

func (db *DatabaseClient) writeFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseClient) GetWhatComicsAreInstalled() map[int]bool {
	mapa := make(map[int]bool, len(db.comicsInfo))
	for i := range db.comicsInfo {
		mapa[i] = true
	}
	return mapa
}

func (db *DatabaseClient) FindComicsByStringUsingIndex(queryWords []globalModel.StemmedWord) []globalModel.ComicInfoToOtput {
	dbqueryWords := convertor.FromGlobalToDBKeywords(queryWords)

	ComicsesWithQueryedWords := make([][]dbModel.OutputProcessModel, len(dbqueryWords))

	for i, queryWord := range dbqueryWords {
		for _, comic := range db.indexTable[queryWord.Word] {
			ComicsesWithQueryedWords[i] = append(ComicsesWithQueryedWords[i], dbModel.OutputProcessModel{
				Num:    comic.ComicsID,
				Weight: comic.Weight,
				Word:   queryWord.Word,
			})
		}
	}

	rankedComics := db.ranking(queryWords, ComicsesWithQueryedWords)

	res := make([]globalModel.ComicInfoToOtput, 0, len(rankedComics))
	for _, comic := range rankedComics {
		res = append(res, globalModel.ComicInfoToOtput{
			Num: comic.Num,
			Img: db.comicsInfo[comic.Num].URL,
		})
	}
	return res
}

func (db *DatabaseClient) FindComicsByStringNotUsingIndex(queryWords []globalModel.StemmedWord) []globalModel.ComicInfoToOtput {
	dbqueryWords := convertor.FromGlobalToDBKeywords(queryWords)

	comicsWithQueryedWords := make([][]dbModel.OutputProcessModel, len(dbqueryWords))

	for i, word := range dbqueryWords {
		for j, comic := range db.comicsInfo {
			for _, comicWord := range comic.Keywords {
				if word.Word == comicWord.Word {
					comicsWithQueryedWords[i] = append(comicsWithQueryedWords[i], dbModel.OutputProcessModel{
						Num:    j,
						Weight: comicWord.Count,
						Word:   word.Word,
					})
				}
			}
		}
	}

	rankedComics := db.ranking(queryWords, comicsWithQueryedWords)

	res := make([]globalModel.ComicInfoToOtput, 0, len(rankedComics))
	for _, comic := range rankedComics {
		res = append(res, globalModel.ComicInfoToOtput{
			Num: comic.Num,
			Img: db.comicsInfo[comic.Num].URL,
		})
	}
	return res
}

func (db *DatabaseClient) ranking(queryWords []globalModel.StemmedWord, suitableWords [][]dbModel.OutputProcessModel) []dbModel.OutputProcessModel {
	var res []dbModel.OutputProcessModel

	for _, some := range suitableWords {
		sort.Slice(some, func(i, j int) bool {
			return some[i].Weight > some[j].Weight
		})
	}

	master := suitableWords[0]
	bitmap := make(map[int]bool)
	for i := 1; i < len(suitableWords); i++ {
		for j := 0; j < len(suitableWords[i]); j++ {
			for _, word := range master {
				if suitableWords[i][j].Num == word.Num && !bitmap[suitableWords[i][j].Num] {
					res = append(res, suitableWords[i][j])
					bitmap[suitableWords[i][j].Num] = true
				}
			}
		}
		master = res
	}

	sort.Slice(queryWords, func(i, j int) bool {
		return queryWords[i].Count > queryWords[j].Count
	})

	for i := 0; i < len(queryWords); i++ {
		mostRelevantWord := i
		for j := 0; j < len(suitableWords[mostRelevantWord]); j++ {
			if _, ok := bitmap[suitableWords[mostRelevantWord][j].Num]; ok {
				continue
			}
			res = append(res, suitableWords[mostRelevantWord][j])
			if len(res) == 10 {
				return res
			}
		}
	}

	return res
}
