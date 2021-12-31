package service

type ItemMgr struct {
	GenID int           `json:"gen_id"`
	Items map[int]*Item `json:"items"`
}

type Item struct {
	ID        int    `json:"id"`
	Desc      string `json:"desc"`
	ClusterID int    `json:"cluster_id"`
	Node      string `json:"node"`
	User      string `json:"user"`
	CreateAt  int64  `json:"create_at"`
	UpdateAt  int64  `json:"update_at"`
	Online    bool   `json:"online"`
}
