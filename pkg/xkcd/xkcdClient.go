package xkcdClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Alphonnse/yaxkcdro/pkg/models"
	"github.com/Alphonnse/yaxkcdro/pkg/xkcd/convertor"
	xkcdModel "github.com/Alphonnse/yaxkcdro/pkg/xkcd/models"
)

type XkcdClient struct {
	resourceURL string
	client      http.Client
}

func NewXkcdClient(resourceURL string) *XkcdClient {
	return &XkcdClient{
		resourceURL: resourceURL,
		client:      http.Client{Timeout: 2 * time.Second},
	}
}

func (c *XkcdClient) GetComicsFromResource(lastDownloadedComic int) (*models.ComicInfoGlobal, error) {
	comicNumber := lastDownloadedComic + 1

	resp, err := c.client.Get(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, comicNumber))
	if err != nil {
		return nil, fmt.Errorf("Error getting comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("No comics with ID %d on %s", comicNumber, c.resourceURL)
	}

	var comicInfo xkcdModel.ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comicInfo)
	if err != nil {
		return nil, fmt.Errorf("Error parsing comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}

	return convertor.FromXkcdClientToGlobal(comicInfo), nil
}

func (c *XkcdClient) GetComicsCountOnResource() (int, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/info.0.json", c.resourceURL))
	if err != nil {
		return 0, fmt.Errorf("Error getting comics count on %s: %s", c.resourceURL, err.Error())
	}
	defer resp.Body.Close()

	var comicInfo xkcdModel.ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comicInfo)
	if err != nil {
		return 0, fmt.Errorf("Error parsing comic from %s: %s", c.resourceURL, err.Error())
	}

	return comicInfo.Num, nil
}
