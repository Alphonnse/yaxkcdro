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
		return nil, fmt.Errorf("Error getting comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got status %d for comics %d", resp.StatusCode, comicNumber)
	}

	var comicInfo xkcdModel.ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comicInfo)
	if err != nil {
		return nil, fmt.Errorf("Error parsing comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}

	return convertor.FromXkcdClientToGlobal(comicInfo), nil
}

func (c *XkcdClient) GetComicsCountOnResource() (int, error) {
	count, err := comicsCountTemp(false, 0)
	if err != nil {
		return 0, fmt.Errorf("Error getting comics count: %s", err.Error())
	}

	if count > 0 {
		return count, nil
	}

	i := 1
	indent := 1000
	var resp *http.Response
	for {
		resp, err = c.client.Get(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, i))
		if err != nil {
			return 0, fmt.Errorf("Error getting comics count on %s: %s", c.resourceURL, err.Error())
		}
		defer resp.Body.Close() 

		if resp.StatusCode == http.StatusNotFound {
			i -= indent
			indent /= 2
		} else {
			i += indent
		}

		if indent == 0 {
			break
		}
	}

	count = i

	if count > 0 {
		_, err = comicsCountTemp(true, count)
		if err != nil {
			return 0, fmt.Errorf("Error writing comics count: %s", err.Error())
		}
	}

	return count, nil
}

type ComicsCountData struct {
	LastRequest time.Time `yaml:"lastRequest"`
	Count       int       `yaml:"count"`
}

func comicsCountTemp(shouldUpdate bool, newCount int) (int, error) {
	path := ".tmp.yaml"
	var data ComicsCountData

	if !shouldUpdate {
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

	// Update the count data and write it to the file.
	data = ComicsCountData{
		LastRequest: time.Now(),
		Count:       newCount,
	}

	updatedData, err := yaml.Marshal(&data)
	if err != nil {
		return 0, err
	}

	err = os.WriteFile(path, updatedData, 0644)
	if err != nil {
		return 0, err
	}

	return data.Count, nil
}
