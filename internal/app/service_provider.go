package app

import (
	"log"

	"github.com/Alphonnse/yaxkcdro/internal/UseCase/api"
	"github.com/Alphonnse/yaxkcdro/internal/UseCase/repository"
	"github.com/Alphonnse/yaxkcdro/internal/UseCase/stemmer"
	"github.com/Alphonnse/yaxkcdro/internal/UseCase/xkcd"
	"github.com/Alphonnse/yaxkcdro/internal/config"
)

type serviceProvider struct {
	pathConfig string

	apiUseCaseService api.APIUseCaseService
	configService     config.ConfigService
	xkcdService       xkcd.XkcdService
	repositoryService repository.RepositoryService
	stemmerService    stemmer.StemmerService
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

func (s *serviceProvider) StemmerService() stemmer.StemmerService {
	if s.stemmerService == nil {
		stemmer := stemmer.NewStemmerClient(s.configService.GetPathStopwords()) 
		s.stemmerService = stemmer
	}
	return s.stemmerService
}

func (s *serviceProvider) XkcdService() xkcd.XkcdService {
	if s.xkcdService == nil {
		xkcd := xkcd.NewXkcdClient(s.configService.GetResourceURL())
		s.xkcdService = xkcd
	}
	return s.xkcdService
}

func (s *serviceProvider) RepositoryService() repository.RepositoryService {
	if s.repositoryService == nil {
		repository, err := repository.NewRepository(s.configService.GetPathDB(), s.configService.GetPathIndex())
		if err != nil {
			log.Fatalf("Field creating database client: %s", err.Error())
		}
		s.repositoryService = repository
	}
	return s.repositoryService
}

func (s *serviceProvider) APIUseCaseService() api.APIUseCaseService {
	if s.apiUseCaseService == nil {
		apiUseCase := api.NewApi(s.RepositoryService())
		s.apiUseCaseService = apiUseCase
	}
	return s.apiUseCaseService
}

// func (s *serviceProvider) DatabaseService() database.DatabaseService {
// 	if s.databaseService == nil {
// 		database, err := database.NewDatabaseClient(s.configService.GetPathDB(), s.configService.GetPathIndex())
// 		if err != nil {
// 			log.Fatalf("Field creating database client: %s", err.Error())
// 		}
// 		s.databaseService = database
// 	}
// 	return s.databaseService
// }
