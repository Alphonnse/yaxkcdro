package xkcd

import (
	"fmt"

	xkcdadapter "github.com/Alphonnse/yaxkcdro/internal/adapters/xkcdAdapter"
	"github.com/Alphonnse/yaxkcdro/internal/models"
	"github.com/Alphonnse/yaxkcdro/pkg/xkcd"
)

type XkcdClient struct {
	xkcdClient *xkcd.Xkcd
}

func NewXkcdClient(resurceURL string) *XkcdClient {
	return &XkcdClient{
		xkcdClient: xkcd.NewXkcdClient(resurceURL),
	}
}

func (c *XkcdClient) GetComicsCountOnResource() (int, error) {
	return c.xkcdClient.GetComicsCount()
}

func (c *XkcdClient) GetComicsFromResource(comicsNumber int) (*models.ComicsInfoUC, error) {
	comicsInfo, err := c.xkcdClient.GetComics(comicsNumber)
	if err != nil {
		return nil, fmt.Errorf("Field to install comics: %s", err.Error())
	}

	return xkcdadapter.FromXkcdGetComicsToGlobal(comicsInfo), nil
}
