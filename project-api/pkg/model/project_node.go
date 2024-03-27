package model

type ProjectNodeTree struct {
	Id       int64              `json:"id"`
	Node     string             `json:"node"`
	Title    string             `json:"title"`
	IsMenu   int                `json:"is_menu"`
	IsLogin  int                `json:"is_login"`
	IsAuth   int                `json:"is_auth"`
	Pnode    string             `json:"pnode"`
	Children []*ProjectNodeTree `json:"children"`
}

type ProjectNodeAuthTree struct {
	Id       int64                  `json:"id"`
	Node     string                 `json:"node"`
	Title    string                 `json:"title"`
	IsMenu   int                    `json:"is_menu"`
	IsLogin  int                    `json:"is_login"`
	IsAuth   int                    `json:"is_auth"`
	Pnode    string                 `json:"pnode"`
	Key      string                 `json:"key"`
	Checked  bool                   `json:"checked"`
	Children []*ProjectNodeAuthTree `json:"children"`
}
