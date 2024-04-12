package convertor

import (
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)


func FromGlobalToDBKeywords(globalMap []globalModel.Keyword) []dbModel.Keyword {
	dbKeywords := make([]dbModel.Keyword, 0, len(globalMap))
	for _, keyword := range globalMap {
		dbKeywords = append(dbKeywords, dbModel.Keyword{
			Word:  keyword.Word,
			Count: keyword.Count,
		})
	}
	return dbKeywords
}

func FromGlobalToDBComicsInfo(globalMap globalModel.ComicInfoGlobal) *dbModel.DBComicsInfo {
	return &dbModel.DBComicsInfo{
		URL:      globalMap.Img,
		Keywords: FromGlobalToDBKeywords(globalMap.Keywords),
	}
}
