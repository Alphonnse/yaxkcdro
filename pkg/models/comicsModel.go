package models

type ComicInfoGlobal struct {
	Num        int           `json:"num"`
	Transcript string        `json:"transcript"`
	Alt        string        `json:"alt"`
	Keywords   []StemmedWord `json:"words"`
	Img        string        `json:"img"`
}

type ComicInfoToOtput struct {
	Num      int      `json:"num"`
	Img      string   `json:"img"`
}

type StemmedWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}
