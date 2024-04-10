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
		client:      http.Client{Timeout: 8 * time.Second},
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
	count, needToRewrite, err := comicsCountTemp(false, 0)
	if err != nil {
		return 0, fmt.Errorf("Error getting comics count, %s", err.Error())
	}
	if needToRewrite {
		i := 1
		indent := 1000
		prev := 0
		var resp *http.Response
		for {
			if prev == i {
				break
			}
			resp, err = c.client.Get(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, i))
			if err != nil {
				return 0, fmt.Errorf("Error getting comics count on %s: %s", c.resourceURL, err.Error())
			}

			prev = i
			if resp.StatusCode == http.StatusNotFound {
				i -= indent
				indent /= 2
			}

			i += indent
		}
		if resp != nil {
			resp.Body.Close()
		}

		count = i

		_, _, err = comicsCountTemp(true, count)
		if err != nil {
			return 0, fmt.Errorf("Error writing comics count, %s", err.Error())
		}
	}
	return count, nil
}

type countInTmp struct {
	LastRequest time.Time `yaml:"lastRequest"`
	Count       int       `yaml:"count"`
}

func comicsCountTemp(shouldBeRewrited bool, newCount int) (CountToRewrite int, shouldBeRewriten bool, err error) {
	path := ".tmp.yaml"
	var count int

	if !shouldBeRewrited {
		_, err := os.Stat(path)
		if err != nil {
			return 0, true, nil
		}
		if os.IsNotExist(err) {
			return 0, true, nil
		}

		configFile, err := os.ReadFile(path)
		if err != nil {
			return 0, false, err
		}

		var readenConunt countInTmp
		err = yaml.Unmarshal(configFile, &readenConunt)
		if err != nil {
			return 0, false, err
		}
		if time.Now().Sub(readenConunt.LastRequest) > time.Hour*4 {
			return 0, true, nil
		}
		count = readenConunt.Count
	} else {
		countToWrite := countInTmp{
			LastRequest: time.Now(),
			Count:       newCount,
		}
		updated, err := yaml.Marshal(&countToWrite)
		if err != nil {
			return 0, false, err
		}
		err = os.WriteFile(path, updated, 0644)
		if err != nil {
			return 0, false, err
		}
	}
	return count, false, nil
}
