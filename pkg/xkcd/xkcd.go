package xkcd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Xkcd struct {
	resourceURL string
	client      http.Client
}

func NewXkcdClient(resourceURL string) *Xkcd {
	return &Xkcd{
		resourceURL: resourceURL,
		client:      http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Xkcd) GetComics(comicNumber int) (*ComicsInfo, error) {

	resp, err := c.client.Get(fmt.Sprintf("%s/%d/info.0.json", c.resourceURL, comicNumber))
	if err != nil {
		return nil, fmt.Errorf("can not get comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %d for comics %d", resp.StatusCode, comicNumber)
	}

	var comicsInfo ComicsInfo
	err = json.NewDecoder(resp.Body).Decode(&comicsInfo)
	if err != nil {
		return nil, fmt.Errorf("can not parse comic %d from %s: %s", comicNumber, c.resourceURL, err.Error())
	}

	return &comicsInfo, nil
}

func (c *Xkcd) GetComicsCount() (int, error) {
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

		if resp.StatusCode == http.StatusNotFound && attempt < 2 {
			attempt++
			count += indent
		} else if resp.StatusCode == http.StatusNotFound && attempt >= 2 {
			count -= 3 * indent
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

const path = ".tmp.yaml"

func readComicsCount() (int, error) {
	var data comicsCountData

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
	data := comicsCountData{
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
