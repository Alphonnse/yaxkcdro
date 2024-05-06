package repoadapter

import (
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/database"
)

func FromDatabaseGetSavedComicsToGlobal(comics map[int]database.DBStorageComicsInfo) map[int]models.ComicsInfoUC {
	out := make(map[int]models.ComicsInfoUC, len(comics))

	for id, comic := range comics {
		out[id] = models.ComicsInfoUC{
			Num:      id,
			Keywords: fromDatabaseKeywordsToGlobal(comic.Keywords),
			Img:      comic.URL,
		}
	}

	return out
}

func fromDatabaseKeywordsToGlobal(keywords []database.Keyword) []models.StemmedWordUC {
	out := make([]models.StemmedWordUC, 0, len(keywords))

	for _, keyword := range keywords {
		out = append(out, models.StemmedWordUC{
			Word:  keyword.Word,
			Count: keyword.Count,
		})
	}

	return out
}

func FromUCSaveComicsToDatabase(comics models.ComicsInfoUC) database.DBComicsInfoToInOut {
	return database.DBComicsInfoToInOut{
		Num:      comics.Num,
		Img:      comics.Img,
		Keywords: fromUCKeywordsToDatabase(comics.Keywords),
	}
}

func fromUCKeywordsToDatabase(keywords []models.StemmedWordUC) []database.Keyword {
	out := make([]database.Keyword, 0, len(keywords))

	for _, keyword := range keywords {
		out = append(out, database.Keyword{
			Word:  keyword.Word,
			Count: keyword.Count,
		})
	}

	return out
}
