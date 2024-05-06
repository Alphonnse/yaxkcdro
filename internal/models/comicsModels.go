package models

type ComicsInfoUC struct {
	Num        int             `json:"num"`
	SafeTitle  string          `json:"safe_title"`
	Transcript string          `json:"transcript"`
	Alt        string          `json:"alt"`
	Keywords   []StemmedWordUC `json:"words"`
	Img        string          `json:"img"`
}

type StemmedWordUC struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}
