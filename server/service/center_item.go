package service

type ItemMgr struct {
	GenID int           `json:"gen_id"`
	Items map[int]*Item `json:"items"`
}

type Item struct {
	ID         int    `json:"id"`
	CreateAt   int64  `json:"create_at"`
	CreateName string `json:"create_name"`
	NodeName   string `json:"node_name"`
	UpdateAt   int64  `json:"update_at"`
	Online     bool   `json:"online"`
}

var (
	itemMgr *ItemMgr
)
