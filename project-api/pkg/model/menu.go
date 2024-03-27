package model

type Menu struct {
	Id         int64   `json:"id"`
	Pid        int64   `json:"pid"`
	Title      string  `json:"title"`
	Icon       string  `json:"icon"`
	Url        string  `json:"url"`
	FilePath   string  `json:"file_path"`
	Params     string  `json:"params"`
	Node       string  `json:"node"`
	Sort       int32   `json:"sort"`
	Status     int32   `json:"status"`
	CreateBy   int64   `json:"create_by"`
	IsInner    int32   `json:"is_inner"`
	Values     string  `json:"values"`
	ShowSlider int32   `json:"show_slider"`
	StatusText string  `json:"statusText"`
	InnerText  string  `json:"innerText"`
	FullUrl    string  `json:"fullUrl"`
	Children   []*Menu `json:"children"`
}
