package app

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Alphonnse/yaxkcdro/pkg/database"
	"github.com/Alphonnse/yaxkcdro/pkg/words"
	xkcdClient "github.com/Alphonnse/yaxkcdro/pkg/xkcd"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp() (*App, error) {
	a := &App{}

	a.InitDeps("config.yaml")

	return a, nil
}

func (a *App) InitDeps(pathConfig string) {
	a.InitServiceProvider(pathConfig)
	a.InitApp()
}

func (a *App) InitServiceProvider(pathConfig string) {
	a.serviceProvider = newServiceProvider(pathConfig)
}

func (a *App) InitApp() {
	a.serviceProvider.ConfigService()
	a.serviceProvider.XkcdService()
	a.serviceProvider.DatabaseService()
	a.serviceProvider.StemmerService()
}

func (a *App) RunApp() error {
	err := downloadComics(a.serviceProvider.stemmerService, a.serviceProvider.xkcdService, a.serviceProvider.databaseService)
	if err != nil {
		return fmt.Errorf("Error downloading comics: %s\n", err.Error())
	}

	comicsID, comicsCount, err := readArgs()
	if err != nil {
		flag.Usage()
		return fmt.Errorf("Error reading args: %s\n", err.Error())
	}

	comicses := a.serviceProvider.databaseService.GetComicsInfo(comicsID, comicsCount)
	for _, comic := range comicses {
		fmt.Printf("\nID: %d\nDescription: %s\nImg: %s\n\n", comic.Num, comic.Keywords, comic.Img)
	}
	return nil
}

func readArgs() (int, int, error) {
	var comicsID int
	var comicsCount int
	flag.IntVar(&comicsID, "o", 0, "First comic ID")
	flag.IntVar(&comicsCount, "n", 0, "Count of comics")
	flag.Parse()
	if comicsID <= 0 && comicsCount <= 0 {
		return 0, 0, errors.New("-o and -n should be positive")
	}

	return comicsID, comicsCount, nil
}

func downloadComics(stemmer words.StemmerService, xkcd xkcdClient.XkcsService, database database.DatabaseService) error {
	comicsCountOnResource, err := xkcd.GetComicsCountOnResource()
	if err != nil {
		return fmt.Errorf("Error counting comics on: %s", err.Error())
	}
	lastDownloadedComic, err := database.FindLastDownloadedComic()
	if err != nil {
		return fmt.Errorf("Error finding last downloaded comic: %s", err.Error())
	}

	if comicsCountOnResource == lastDownloadedComic {
		fmt.Println("All comics already downloaded")
		return nil
	}

	for comicsCountOnResource > lastDownloadedComic {
		currentComic := lastDownloadedComic + 1
		fmt.Printf("\rDownloading comic %d/%d", currentComic, comicsCountOnResource)

		comicsInfo, err := xkcd.GetComicsFromResource(lastDownloadedComic)
		if err != nil {
			fmt.Println()
			log.Printf("Error getting comic %d: %s", currentComic, err.Error())
			lastDownloadedComic++
			continue
		}

		comicsInfo, err = stemmer.Stem(*comicsInfo)
		if err != nil {
			fmt.Println()
			log.Printf("Error stemming comic %d: %s", currentComic, err.Error())
			lastDownloadedComic++
			continue
		}

		err = database.InsertComicsIntoDB(*comicsInfo)
		if err != nil {
			fmt.Println()
			log.Printf("Error inserting comic %d into database: %s", currentComic, err.Error())
			continue
		}
		lastDownloadedComic++
	}

	fmt.Println("\nAll comics are downloaded")
	return nil
}
