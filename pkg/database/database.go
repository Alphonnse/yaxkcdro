package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

type Database struct {
	pathToDBFile    string
	pathToIndexFile string

	comicsInfo map[int]DBStorageComicsInfo
	indexTable map[string][]IndexModel

	mu sync.RWMutex
}

func NewDatabase(pathToDBFile, pathToIndexFile string) (*Database, error) {
	client := &Database{
		pathToDBFile:    pathToDBFile,
		comicsInfo:      make(map[int]DBStorageComicsInfo),
		pathToIndexFile: pathToIndexFile,
		indexTable:      make(map[string][]IndexModel),
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

func (db *Database) readFile(index bool, pathToFile string) error {
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

func (db *Database) GetComics() (map[int]DBStorageComicsInfo, error) {
	return db.comicsInfo, nil
}

func (db *Database) SaveComics(comicsInfo DBComicsInfoToInOut) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.comicsInfo[comicsInfo.Num] = DBStorageComicsInfo{
		URL:      comicsInfo.Img,
		Keywords: comicsInfo.Keywords,
	}
	db.processIndex(&comicsInfo)

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

	return nil
}

func (db *Database) processIndex(comics *DBComicsInfoToInOut) {
	for _, comicsWord := range comics.Keywords {
		db.indexTable[comicsWord.Word] = append(db.indexTable[comicsWord.Word], IndexModel{
			ComicsID: comics.Num,
			Weight:   comicsWord.Count,
		})
	}
}

func (db *Database) writeFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Find(queryWords []Keyword) []DBComicsInfoToInOut {
	ComicsesWithQueryedWords := make([][]findOutputProcessModel, len(queryWords))

	for i, queryWord := range queryWords {
		for _, comic := range db.indexTable[queryWord.Word] {
			ComicsesWithQueryedWords[i] = append(ComicsesWithQueryedWords[i], findOutputProcessModel{
				Num:    comic.ComicsID,
				Weight: comic.Weight,
				Word:   queryWord.Word,
			})
		}
	}

	rankedComics := db.rank(queryWords, ComicsesWithQueryedWords)

	res := make([]DBComicsInfoToInOut, 0, len(rankedComics))
	for _, comic := range rankedComics {
		res = append(res, DBComicsInfoToInOut{
			Num: comic.Num,
			Img: db.comicsInfo[comic.Num].URL,
		})
	}
	return res
}

func (db *Database) rank(queryWords []Keyword, suitableWords [][]findOutputProcessModel) []findOutputProcessModel {
	var res []findOutputProcessModel

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
