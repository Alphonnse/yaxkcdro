package xkcdModel

type ComicInfo struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	SafeTitle  string `json:"safe_title"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}
