package server

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

type nodeStore struct {
	storeBase
}

func (store *nodeStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&nodes, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
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

type userMgrStore struct {
	storeBase
}

func (store *userMgrStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&userMgr, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			userMgr = &UserMgr{
				Admin:   nil,
				UserMap: map[string]*User{},
			}
		}
		return
	}
	return
}

func (store *userMgrStore) Save() {
	_ = util.EncodeJsonToFile(userMgr, store.filename)
}

type templateStore struct {
	storeBase
}

func (store *templateStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&temps, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	return
}

func (store *templateStore) Save() {
	_ = util.EncodeJsonToFile(temps, store.filename)
}

type cluMgrStore struct {
	storeBase
}

func (store *cluMgrStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&cluMgr, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			cluMgr = &ClusterMgr{
				GenID:    0,
				Clusters: map[int]*Cluster{},
			}
		}
		return
	}
	return
}

func (store *cluMgrStore) Save() {
	_ = util.EncodeJsonToFile(cluMgr, store.filename)
}

type cmdMgrStore struct {
	storeBase
}

func (store *cmdMgrStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&cmdMgr, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			cmdMgr = &CmdMgr{
				CmdMap:  map[string]*Cmd{},
				CmdLogs: map[string][]*CmdLog{},
			}
		}
		return
	}
	return
}

func (store *cmdMgrStore) Save() {
	_ = util.EncodeJsonToFile(cmdMgr, store.filename)
}

type storeName string

const (
	snNode     storeName = "node"
	snUserMgr  storeName = "user_mgr"
	snCmdMgr   storeName = "cmd_mgr"
	snItemMgr  storeName = "item_mgr"
	snTemplate storeName = "template"
	snCluMgr   storeName = "clu_mgr"
)

var (
	userMgr *UserMgr
	itemMgr *ItemMgr
	nodes   = map[string]*Node{}
	temps   = map[string]*Template{} // web
	cluMgr  *ClusterMgr
	cmdMgr  *CmdMgr
)

func init() {
	stores[snNode] = &nodeStore{storeBase{
		file: "node.json",
	}}
	stores[snItemMgr] = &itemMgrStore{storeBase{
		file: "item_mgr.json",
	}}
	stores[snUserMgr] = &userMgrStore{storeBase{
		file: "user_mgr.json",
	}}
	stores[snTemplate] = &templateStore{storeBase{
		file: "template.json",
	}}
	stores[snCluMgr] = &cluMgrStore{storeBase{
		file: "clu_mgr.json",
	}}
	stores[snCmdMgr] = &cmdMgrStore{storeBase{
		file: "cmd_mgr.json",
	}}
}
