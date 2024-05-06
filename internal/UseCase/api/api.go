package api

import (
	"fmt"
	"log"

	"github.com/Alphonnse/yaxkcdro/internal/UseCase/repository"
	"github.com/Alphonnse/yaxkcdro/internal/UseCase/stemmer"
	"github.com/Alphonnse/yaxkcdro/internal/UseCase/xkcd"
)

type API struct {
	repository repository.RepositoryService
	xkcdClient xkcd.XkcdService
	stemmerClient stemmer.StemmerService
}

func NewApi(repo repository.RepositoryService) *API {
	return &API{
		repository: repo,
	}
}

func (a *API) UpdateComicsCount() error {
	comicsCountOnResource, err := a.xkcdClient.GetComicsCountOnResource()
	if err != nil {
		return fmt.Errorf("Error counting comics on: %s", err.Error())
	}
	// log.Printf("Count of comics on resource: %d\n", comicsCountOnResource)
	comicsInDB, err := a.repository.GetSavedComics()
	comicsCountInDB := len(comicsInDB)

	// log.Printf("Count of installed comics in DB (ex.404): %d\n", comicsCountInDB)

	if comicsCountOnResource-1 == comicsCountInDB {
		// тут тоже по факту должен быть log
		// fmt.Println("All comics already downloaded")
		return nil
	}

	// take care about the interruption

	var comicsToInstall []Task
	for i := 1; i <= comicsCountOnResource; i++ {
		if i == 404 {
			continue
		}
		if _, ok := comicsInDB[i]; !ok {
			comicsToInstall = append(comicsToInstall, &comicsInstallerTask{
				comicsID:   i,
				apiUseCase: a,
			})
		}
	}

	return nil
}

func (t *comicsInstallerTask) process() {
	// количество можно будет контролить через канал
	comicsInfo, err := t.apiUseCase.xkcdClient.GetComicsFromResource(t.comicsID)
	if err != nil {
		log.Printf("Error getting comics %d: %s", t.comicsID, err.Error())
		return
	}

	keywords, err := t.apiUseCase.stemmerClient.StemComicsDescription(comicsInfo.SafeTitle, comicsInfo.Transcript, comicsInfo.Alt)
	if err != nil {
		log.Printf("Error stemming comic %d: %s", t.comicsID, err.Error())
		return
	}
	comicsInfo.Keywords = keywords

	// сейвить каждые сколько-то
	err = t.apiUseCase.repository.SaveComics(*comicsInfo)
	if err != nil {
		log.Printf("Error inserting comic %d into database: %s", t.comicsID, err.Error())
		return
	}
}

