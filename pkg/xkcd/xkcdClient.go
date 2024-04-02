package xkcdClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetComicsFromResource(resourceURL string, lastDownloadedComic int) (map[int]ComicInfo, error) {
	client := http.Client{Timeout: 2 * time.Second}
	comics := make(map[int]ComicInfo)

	comicsCount, err := getComicsCountOnResource(resourceURL)
	if err != nil {
		return comics, fmt.Errorf("Error getting number of comics from %s: %w", resourceURL, err)
	}

	for comicNumber := lastDownloadedComic+1; comicNumber <= comicsCount; comicNumber++ {
		resp, err := client.Get(fmt.Sprintf("%s/%d/info.0.json", resourceURL, comicNumber))
		// redirecting is errors
		if err != nil {
			return comics, fmt.Errorf("Error getting comic %d from %s: %w", comicNumber, resourceURL, err)
		}

		// parse comic info and add it to the map
		var comicInfo ComicInfo
		err = json.NewDecoder(resp.Body).Decode(&comicInfo)
		if err != nil {
			return comics, fmt.Errorf("Error parsing comic %d from %s: %w", comicNumber, resourceURL, err)
		}
		comics[comicNumber] = comicInfo

		resp.Body.Close()
	}

	fmt.Printf("Number of comics downloaded: %d\n", comicsCount)
	return comics, nil
}

func getComicsCountOnResource(resourceURL string) (int, error) {
	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(fmt.Sprintf("%s/info.0.json", resourceURL))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var comicInfo ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comicInfo)
	if err != nil {
		return 0, err
	}

	return comicInfo.Num, nil
}
