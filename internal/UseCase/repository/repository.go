package repository

import (
	"fmt"

	repoadapter "github.com/Alphonnse/yaxkcdro/internal/adapters/repoAdapter"
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/database"
)

type Reposiroty struct {
	database *database.Database
}

func NewRepository(pathToDBFile string, pathToIndexFile string) (*Reposiroty, error) {
	database, err := database.NewDatabase(pathToDBFile, pathToIndexFile)
	if err != nil {
		return nil, fmt.Errorf("Field creating database client: %s", err.Error())
	}
	return &Reposiroty{
		database: database,
	}, nil
}

func (r *Reposiroty) GetSavedComics() (map[int]bool, error) {
	SavedComics, err := r.database.GetComics()
	if err != nil {
		return nil, fmt.Errorf("Field getting installed comics: %s", err.Error())
	}
	SavedComicsUC := repoadapter.FromDatabaseGetSavedComicsToGlobal(SavedComics)

	savedComicsList := make(map[int]bool, len(SavedComicsUC))
	for i := range SavedComicsUC {
		savedComicsList[i] = true
	}

	return savedComicsList, nil
}

func (r *Reposiroty) SaveComics(comics models.ComicsInfoUC) error {
	err := r.database.SaveComics(repoadapter.FromUCSaveComicsToDatabase(comics))
	if err != nil {
		return fmt.Errorf("Field saving installed comic: %s", err.Error())
	}
	return nil
}
