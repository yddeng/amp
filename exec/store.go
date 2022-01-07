package exec

import (
	"initial-server/util"
	"os"
	"path"
)

type Store interface {
	Load(dataPath string) error
	Save()
}

type storeBase struct {
	file     string
	filename string
}

var stores = map[storeName]Store{}

func loadStore(dataPath string) (err error) {
	for _, store := range stores {
		if err = store.Load(dataPath); err != nil {
			return
		}
	}
	return
}

func saveStore(names ...storeName) {
	if len(names) == 0 {
		for _, store := range stores {
			store.Save()
		}
	} else {
		for _, name := range names {
			if store, ok := stores[name]; ok {
				store.Save()
			}
		}
	}
}

type appStore struct {
	storeBase
}

func (store *appStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&apps, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	return
}

func (store *appStore) Save() {
	_ = util.EncodeJsonToFile(apps, store.filename)
}

type storeName string

const (
	snApp storeName = "application"
)

var (
	apps = map[int]*Application{}
)

func init() {
	stores[snApp] = &appStore{storeBase{
		file: "application.json",
	}}
}
