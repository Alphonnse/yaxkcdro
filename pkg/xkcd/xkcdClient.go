package xkcdClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Alphonnse/yaxkcdro/pkg/models"
	"github.com/Alphonnse/yaxkcdro/pkg/xkcd/convertor"
	xkcdModel "github.com/Alphonnse/yaxkcdro/pkg/xkcd/models"
	"gopkg.in/yaml.v3"
)

type XkcdClient struct {
	resourceURL string
	client      http.Client
}

func NewXkcdClient(resourceURL string) *XkcdClient {
	return &XkcdClient{
		resourceURL: resourceURL,
		client:      http.Client{Timeout: 5 * time.Second},
	}
}

func (c *XkcdClient) GetComicsFromResource(comicNumber int) (*models.ComicInfoGlobal, error) {

	resp, err := c.client.Get(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, comicNumber))
	if err != nil {
		return nil, fmt.Errorf("can not get comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %d for comics %d", resp.StatusCode, comicNumber)
	}

	var comicInfo xkcdModel.ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comicInfo)
	if err != nil {
		return nil, fmt.Errorf("can not parse comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}

	return convertor.FromXkcdClientToGlobal(comicInfo), nil
}

func (c *XkcdClient) GetComicsCountOnResource() (int, error) {
	count, err := readComicsCount()
	if err != nil {
		return 0, fmt.Errorf("can not get comics count: %s", err.Error())
	}

	if count > 0 {
		return count, nil
	}

	count = 1
	indent := 1000
	var resp *http.Response
	var attempt int
	for {
		resp, err = c.client.Head(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, count))
		if err != nil {
			return 0, fmt.Errorf("can not get comics count on %s: %s", c.resourceURL, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound && attempt < 2{
			attempt++
			count += indent
		} else if resp.StatusCode == http.StatusNotFound && attempt >= 2 {
			count -= 3*indent
			indent /= 2
			attempt = 0
		} else {
			count += indent
		}

		if indent == 0 {
			break
		}
	}

	if count > 0 {
		err = writeComicsCount(count)
		if err != nil {
			return 0, fmt.Errorf("can not write comics count: %s", err.Error())
		}
	}

	return count, nil
}

type ComicsCountData struct {
	LastRequest time.Time `yaml:"lastRequest"`
	Count       int       `yaml:"count"`
}

const path = ".tmp.yaml"

func readComicsCount() (int, error) {
	var data ComicsCountData

	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return 0, nil
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	err = yaml.Unmarshal(configFile, &data)
	if err != nil {
		return 0, err
	}

	if time.Since(data.LastRequest) > time.Hour*4 {
		return 0, nil
	}

	return data.Count, nil
}

func writeComicsCount(newCount int) error {
	var data ComicsCountData

	data = ComicsCountData{
		LastRequest: time.Now(),
		Count:       newCount,
	}

	updatedData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
