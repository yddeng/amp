package center

import (
	"os"
	"path"
)

const (
	rePath   = "center"
	itemFile = "item.json"
)

var (
	dataPath string
)

func LoadData(root string) (err error) {
	dataPath = path.Join(root, rePath)
	_ = os.MkdirAll(dataPath, os.ModePerm)

	return
}
