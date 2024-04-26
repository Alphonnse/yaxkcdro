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

	configPath := "config/config.yaml"
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

	searchString, searchByIndex := readArgs()

	stemmedSearchString, err := a.serviceProvider.stemmerService.StemQueryText(searchString)
	if err != nil {
		return fmt.Errorf("Error stemming text: %s\n", err.Error())
	}

	err = downloadComics(a.serviceProvider)
	if err != nil {
		return fmt.Errorf("Error downloading comics: %s\n", err.Error())
	}

	if searchByIndex {
		for _, coimcs := range a.serviceProvider.databaseService.FindComicsByStringUsingIndex(stemmedSearchString) {
			fmt.Println(coimcs.Img, coimcs.Num)
		}
		return nil
	}
	for _, coimcs := range a.serviceProvider.databaseService.FindComicsByStringNotUsingIndex(stemmedSearchString) {
		fmt.Println(coimcs.Img, coimcs.Num)
	}

	return nil
}

func readArgs() (string, bool) {
	var str string
	cliArgs := os.Args

	if len(cliArgs) == 3 {
		if cliArgs[1] == "-s" {
			return cliArgs[2], false
		} else {
			log.Fatal("Wrong key. Please use -s key to specify a sentence, and -i to search using index ")
		}
	} else if len(cliArgs) == 4 {
		if cliArgs[1] == "-s" && cliArgs[3] == "-i" {
			return cliArgs[2], true
		} else {
			log.Fatal("Wrong key. Please use -s key to specify a sentence, and -i to search using index ")
		}
	} else {
		log.Fatal("Please use -s key to specify a sentence, and -i to search using index")
	}

	return str, false
}

func downloadComics(serviceProvider *serviceProvider) error {
	log.Println("Reading count of comics on resource...")
	comicsCountOnResource, err := serviceProvider.xkcdService.GetComicsCountOnResource()
	if err != nil {
		return fmt.Errorf("Error counting comics on: %s", err.Error())
	}
	log.Printf("Count of comics on resource: %d\n", comicsCountOnResource)

	comicsesInDB := serviceProvider.databaseService.GetWhatComicsAreInstalled()
	comicsCountInDB := len(comicsesInDB)

	log.Printf("Count of installed comics in DB (ex.404): %d\n", comicsCountInDB)

	if comicsCountOnResource-1 == comicsCountInDB {
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

	serviceProvider.databaseService.SetChunkSize(int(float64(comicsCountOnResource-1)*0.05), comicsCountOnResource-1)

	wp := NewWorkerPool(bar, comicsToInstall, serviceProvider, (len(comicsToInstall)/100)+1)
	wp.Run()

	fmt.Println("\nAll comics now downloaded")

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

	keywords, err := t.serviceProvider.stemmerService.StemComicsDesc(comicsInfo.SafeTitle, comicsInfo.Transcript, comicsInfo.Alt)
	if err != nil {
		log.Printf("Error stemming comic %d: %s", t.ComicsID, err.Error())
		return
	}
	comicsInfo.Keywords = keywords

	err = t.serviceProvider.databaseService.InsertComicsIntoFiles(*comicsInfo)
	if err != nil {
		log.Printf("Error inserting comic %d into database: %s", t.ComicsID, err.Error())
		return
	}
	t.bar.Increment()
}
