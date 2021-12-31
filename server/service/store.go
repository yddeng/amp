package service

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

type storeName string

const (
	snNode    storeName = "node"
	snItemMgr storeName = "item_mgr"
	snUser    storeName = "user"
	snAdmin   storeName = "admin"
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

type nodeStore struct {
	storeBase
}

func (store *nodeStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&nodes, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			nodes = map[string]*Node{}
		}
		return
	}
	return
}

func (store *nodeStore) Save() {
	_ = util.EncodeJsonToFile(nodes, store.filename)
}

type itemMgrStore struct {
	storeBase
}

func (store *itemMgrStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&itemMgr, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			itemMgr = &ItemMgr{GenID: 0, Items: map[int]*Item{}}
		}
		return
	}
	return
}

func (store *itemMgrStore) Save() {
	_ = util.EncodeJsonToFile(itemMgr, store.filename)
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
	stores[snNode] = &nodeStore{storeBase{
		file: "node.json",
	}}
	stores[snItemMgr] = &itemMgrStore{storeBase{
		file: "item_mgr.json",
	}}
	stores[snUser] = &storeUser{storeBase{
		file: "user.json",
	}}
	stores[snAdmin] = &adminStore{storeBase{
		file: "admin.json",
	}}
}
