package service

import (
	"initial-sever/util"
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

type storeName string

const (
	snUser  storeName = "user"
	snAdmin storeName = "admin"
)

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

type adminStore struct {
	storeBase
}

func (store *adminStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&admin, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	return
}

func (store *adminStore) Save() {
	_ = util.EncodeJsonToFile(admin, store.filename)
}

type storeUser struct {
	storeBase
}

func (store *storeUser) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&userMap, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			userMap = map[string]*User{}
		}
		return
	}
	return
}

func (store *storeUser) Save() {
	_ = util.EncodeJsonToFile(userMap, store.filename)
}

func init() {
	stores[snUser] = &storeUser{storeBase{
		file: "user.json",
	}}
	stores[snAdmin] = &adminStore{storeBase{
		file: "admin.json",
	}}
}
