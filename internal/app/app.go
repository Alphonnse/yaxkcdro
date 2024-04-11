package app

import (
	"fmt"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp() (*App, error) {
	a := &App{}

	configPath := readArgs()
	a.InitDeps(configPath)

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
	return nil
}

func readArgs() string {
	var str string
	cliArgs := os.Args

	if len(cliArgs) == 3 {
		if cliArgs[1] == "-c" {
			str = cliArgs[2]
		} else {
			log.Fatal("Wrong key. Please use -c key to specify the config file")
		}
	} else {
		log.Fatal("Please use -s key only to specify the config file")
	}
	return str
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

	log.Printf("Count of installed comics in DB (ex.404): %d\n", comicsCountInDB)

	if comicsCountOnResource-1 == comicsCountInDB {
=======
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

	bar := pb.StartNew(comicsCountOnResource - comicsCountInDB)
	bar.Set("prefix", "Downloading comics")
	bar.SetMaxWidth(80)

	var comicsToInstall []Task
	for i := 1; i <= comicsCountOnResource; i++ {
		if i == 404 {
			continue
		}
		if _, ok := comicsesInDB[i]; !ok {
			comicsToInstall = append(comicsToInstall, &ComicsInstallerTask{
				ComicsID:        i,
				serviceProvider: serviceProvider,
				bar:             bar,
			})
		}
	}

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
	ComicsID        int
	serviceProvider *serviceProvider
	bar             *pb.ProgressBar
}

func (t *ComicsInstallerTask) Process() {
	comicsInfo, err := t.serviceProvider.xkcdService.GetComicsFromResource(t.ComicsID)
	if err != nil {
		log.Printf("Error getting comics %d: %s", t.ComicsID, err.Error())
		return
	}

	comicsInfo, err = t.serviceProvider.stemmerService.Stem(*comicsInfo)
	if err != nil {
		log.Printf("Error stemming comic %d: %s", t.ComicsID, err.Error())
		return
	}

	err = t.serviceProvider.databaseService.InsertComicsIntoDB(*comicsInfo)
	if err != nil {
		log.Printf("Error inserting comic %d into database: %s", t.ComicsID, err.Error())
		return
	}
	t.bar.Increment()
}
