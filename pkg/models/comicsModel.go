package models

type ComicInfoGlobal struct {
	Num        int       `json:"num"`
	Transcript string    `json:"transcript"`
	Alt        string    `json:"alt"`
	Keywords   []Keyword `json:"words"`
	Img        string    `json:"img"`
}

type ComicInfoToOtput struct {
	Num      int      `json:"num"`
	Keywords []string `json:"words"`
	Img      string   `json:"img"`
}

type Keyword struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}
