package app

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alphonnse/yaxkcdro/internal/config"
	"github.com/Alphonnse/yaxkcdro/pkg/database"
	"github.com/Alphonnse/yaxkcdro/pkg/words"
	xkcdClient "github.com/Alphonnse/yaxkcdro/pkg/xkcd"
)

type App struct {
	appConfig config.AppConfig
	stemmer   *words.Stemmer
	xkcd      *xkcdClient.XkcdClient
	database  *database.DatabaseClient
}

func NewApp() *App {
	appConfig, err := config.GetAppConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	stemmer := words.NewStemmer(appConfig.PathStopwords)
	xkcd := xkcdClient.NewXkcdClient(appConfig.ResourceURL)
	database, err := database.NewDatabaseClient(appConfig.PathDB)
	if err != nil {
		log.Fatalf("Error creating database client: %v\n", err)
	}

	return &App{
		appConfig: appConfig,
		stemmer:   stemmer,
		xkcd:      xkcd,
		database:  database,
	}
}

func (a *App) RunApp() {
	err := downloadComics(a.stemmer, a.xkcd, a.database)
	if err != nil {
		log.Fatalf("Error downloading comics: %s\n", err.Error())
	}

	comicsID, comicsCount, err := readArgs()
	if err != nil {
		log.Printf("Error reading args: %s\n", err.Error())
		flag.Usage()
		os.Exit(1)
	}

	comicses := a.database.GetComicsInfo(comicsID, comicsCount)
	for _, comic := range comicses {
		fmt.Printf("\nID: %d\nDescription: %s\nImg: %s\n\n", comic.Num, comic.Keywords, comic.Img)
	}
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

func downloadComics(stemmer *words.Stemmer, xkcd *xkcdClient.XkcdClient, database *database.DatabaseClient) error {
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

	// здесь можно вручную установить количество комиксов на сайте, если не нужно
	// закачивать все в БД
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
			log.Printf("\nError stemming comic %d: %s", currentComic, err.Error())
			lastDownloadedComic++
			continue
		}

		err = database.InsertComicsIntoDB(*comicsInfo)
		if err != nil {
			fmt.Println()
			log.Printf("\nError inserting comic %d into database: %s", currentComic, err.Error())
			continue
		}
		lastDownloadedComic++
	}

	fmt.Println("\nAll comics are downloaded")
	return nil
}
