package center

type ItemMgr struct {
	id    int32
	Items map[int32]*Item `json:"items"`
}

type Item struct {
	ItemID int32 `json:"item_id"`
	PID    int32 `json:"pid"`
	Alive  bool  `json:"_"`
}
