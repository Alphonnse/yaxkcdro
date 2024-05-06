package xkcd

import "time"

type comicsCountData struct {
	LastRequest time.Time `yaml:"lastRequest"`
	Count       int       `yaml:"count"`
}

type ComicsInfo struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	SafeTitle  string `json:"safe_title"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}
