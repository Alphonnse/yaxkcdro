package app

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/cheggaaa/pb/v3"
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
	err := downloadComics(a.serviceProvider)
	if err != nil {
		return fmt.Errorf("Error downloading comics: %s\n", err.Error())
	}
	//
	// comicsID, comicsCount, err := readArgs()
	// if err != nil {
	// 	flag.Usage()
	// 	return fmt.Errorf("Error reading args: %s\n", err.Error())
	// }
	//
	// comicses := a.serviceProvider.databaseService.GetComicsInfo(comicsID, comicsCount)
	// for _, comic := range comicses {
	// 	fmt.Printf("\nID: %d\nDescription: %s\nImg: %s\n\n", comic.Num, comic.Keywords, comic.Img)
	// }
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

func downloadComics(serviceProvider *serviceProvider) error {
	log.Println("Reading count of comics on resource...")
	comicsCountOnResource, err := serviceProvider.xkcdService.GetComicsCountOnResource()
	if err != nil {
		return fmt.Errorf("Error counting comics on: %s", err.Error())
	}
	log.Printf("Count of comics on resource: %d\n", comicsCountOnResource)

	comicsesInDB := serviceProvider.databaseService.GetInstalledComics()
	comicsCountInDB := len(comicsesInDB)

	if comicsCountOnResource-1 == comicsCountInDB {
		fmt.Println("All comics already downloaded")
		return nil
	}

	var comicsToInstall []Task
	for i := 1; i <= comicsCountOnResource; i++ {
		if i == 404 {
			continue
		}
		if _, ok := comicsesInDB[i]; !ok {
			comicsToInstall = append(comicsToInstall, &ComicsInstallerTask{ComicsID: i})
		}
	}

	bar := pb.StartNew(len(comicsToInstall))
	bar.Set("prefix", "Downloading comics")
	bar.SetMaxWidth(80)

	// Эта функция необходима, без него будет ошибка в БД
	// почему он не в конструктору? Потому, что в конструкторе я не знаю количество комиксов
	serviceProvider.databaseService.SetChunkSize(int(float64(comicsCountOnResource-1)*0.05), comicsCountOnResource-1)

	// +1 так как если число меньше 100, то при делении получается 0
	wp := NewWorkerPool(bar, comicsToInstall, serviceProvider, (len(comicsToInstall)/100)+1)
	wp.Run()

	fmt.Println("\nAll comics downloaded")

	return nil
}

type ComicsInstallerTask struct {
	ComicsID int
}

func (t *ComicsInstallerTask) Process(serviceProvider *serviceProvider, bar *pb.ProgressBar) {
	comicsInfo, err := serviceProvider.xkcdService.GetComicsFromResource(t.ComicsID)
	if err != nil {
		log.Printf("Error getting comics %d: %s", t.ComicsID, err.Error())
		return
	}

	comicsInfo, err = serviceProvider.stemmerService.Stem(*comicsInfo)
	if err != nil {
		log.Printf("Error stemming comic %d: %s", t.ComicsID, err.Error())
		return
	}

	err = serviceProvider.databaseService.InsertComicsIntoDB(*comicsInfo)
	if err != nil {
		log.Printf("Error inserting comic %d into database: %s", t.ComicsID, err.Error())
		return
	}
	bar.Increment()
}
