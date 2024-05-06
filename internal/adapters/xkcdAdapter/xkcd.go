package xkcdadapter

import (
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/xkcd"
)

func FromXkcdGetComicsToGlobal(comics *xkcd.ComicsInfo) *models.ComicsInfoUC {
	return &models.ComicsInfoUC{
		Num:        comics.Num,
		SafeTitle:  comics.SafeTitle,
		Transcript: comics.Transcript,
		Alt:        comics.Alt,
		Img:        comics.Img,
	}
}
