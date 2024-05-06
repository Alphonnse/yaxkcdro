package dbModel

type DBComicsInfo struct {
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

type OutputProcessModel struct {
	Num    int
	Weight int
	Word   string
}
