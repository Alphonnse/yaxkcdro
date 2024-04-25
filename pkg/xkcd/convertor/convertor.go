package convertor

import (
	"github.com/Alphonnse/yaxkcdro/pkg/models"
	xkcdClient "github.com/Alphonnse/yaxkcdro/pkg/xkcd/models"
)

func FromXkcdClientToGlobal(comicsInfo xkcdClient.ComicInfo) *models.ComicInfoGlobal {
	return &models.ComicInfoGlobal{
		Num:        comicsInfo.Num,
		Transcript: comicsInfo.Transcript,
		SafeTitle:  comicsInfo.SafeTitle,
		Alt:        comicsInfo.Alt,
		Img:        comicsInfo.Img,
	}
}
