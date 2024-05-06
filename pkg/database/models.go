package database

type DBStorageComicsInfo struct {
	URL      string    `json:"url"`
	Keywords []Keyword `json:"words"`
}

type Keyword struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type IndexModel struct {
	ComicsID int `json:"comicsID"`
	Weight   int `json:"weight"`
}

type DBComicsInfoToInOut struct {
	Num      int       `json:"num"`
	Img      string    `json:"url"`
	Keywords []Keyword `json:"words"`
}

type findOutputProcessModel struct {
	Num    int
	Weight int
	Word   string
}

// type OutputComicsInfoModel struct {
// 	Num int
// 	Img string
// }
