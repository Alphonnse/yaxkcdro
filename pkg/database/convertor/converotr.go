package convertor

import (
	dbModel "github.com/Alphonnse/yaxkcdro/pkg/database/models"
	globalModel "github.com/Alphonnse/yaxkcdro/pkg/models"
)

func FromGlobalToDBComicsInfo(globalMap globalModel.ComicInfoGlobal) *dbModel.DBComicsInfo {
	return &dbModel.DBComicsInfo{
		URL:      globalMap.Img,
		Keywords: globalMap.Keywords,
	}
}
