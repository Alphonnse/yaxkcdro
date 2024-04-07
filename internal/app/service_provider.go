package app

import (
	"log"

	"github.com/Alphonnse/yaxkcdro/internal/config"
	"github.com/Alphonnse/yaxkcdro/pkg/database"
	"github.com/Alphonnse/yaxkcdro/pkg/words"
	xkcdClient "github.com/Alphonnse/yaxkcdro/pkg/xkcd"
)

type serviceProvider struct {
	pathConfig string

	configService config.ConfigService
	xkcdService xkcdClient.XkcsService
	databaseService database.DatabaseService
	stemmerService words.StemmerService
}

func newServiceProvider(pathToCOnfig string) *serviceProvider {
	return &serviceProvider{
		pathConfig: pathToCOnfig,
	}
}

func (s *serviceProvider) ConfigService() config.ConfigService {
	if s.configService == nil {
		configService, err := config.NewAppConfig(s.pathConfig)
		if err != nil {
			log.Fatalf("Field to get config: %s", err.Error())
		}
		s.configService = configService
	}
	return s.configService
}

func (s *serviceProvider) StemmerService() words.StemmerService {
	if s.stemmerService == nil {
		stemmer := words.NewStemmer(s.configService.GetPathStopwords())
		s.stemmerService = stemmer
	}
	return s.stemmerService
}

func (s *serviceProvider) XkcdService() xkcdClient.XkcsService {
	if s.xkcdService == nil {
		xkcd := xkcdClient.NewXkcdClient(s.configService.GetResourceURL())
		s.xkcdService = xkcd
	}
	return s.xkcdService
}

func (s *serviceProvider) DatabaseService() database.DatabaseService {
	if s.databaseService == nil {
		database, err := database.NewDatabaseClient(s.configService.GetPathDB())
		if err != nil {
			log.Fatalf("Field creating database client: %s", err.Error())
		}
		s.databaseService = database
	}
	return s.databaseService
}
