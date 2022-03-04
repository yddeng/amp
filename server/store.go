package server

import (
	"amp/util"
	"errors"
	"log"
	"os"
	"path"
)

type Store interface {
	Load(dataPath string) error
	Save() error
}

type storeBase struct {
	file     string
	filename string
}

var (
	stores   = map[storeName]Store{}
	needSave = map[storeName]bool{}
)

func loadStore(dataPath string) (err error) {
	for name, store := range stores {
		if err = store.Load(dataPath); err != nil {
			return errors.New(string(name) + " : " + err.Error())
		}
	}
	return
}

func saveStore(names ...storeName) {
	if len(names) == 0 {
		for name := range stores {
			needSave[name] = true
		}
	} else {
		for _, name := range names {
			needSave[name] = true
		}
	}
}

func doSave(final bool) {
	if final {
		for name, store := range stores {
			if err := store.Save(); err != nil {
				log.Printf("store %s save failed, %s\n", name, err)
			}
		}
	} else {
		for name := range needSave {
			if store, ok := stores[name]; ok {
				if err := store.Save(); err != nil {
					log.Printf("store %s save failed, %s\n", name, err)
				}
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

func (store *nodeStore) Save() error {
	return util.EncodeJsonToFile(nodes, store.filename)
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

func (store *userMgrStore) Save() error {
	return util.EncodeJsonToFile(userMgr, store.filename)
}

type processMgrStore struct {
	storeBase
}

func (store *processMgrStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&processMgr, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
			processMgr = &ProcessMgr{
				GenID:   0,
				Groups:  map[string]struct{}{},
				Process: map[int]*Process{},
			}
		}
		return
	}
	return
}

func (store *processMgrStore) Save() error {
	return util.EncodeJsonToFile(processMgr, store.filename)
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
				CmdMap:  map[int]*Cmd{},
				CmdLogs: map[int][]*CmdLog{},
			}
		}
		return
	}
	return
}

func (store *cmdMgrStore) Save() error {
	return util.EncodeJsonToFile(cmdMgr, store.filename)
}

type kvStore struct {
	storeBase
}

func (store *kvStore) Load(dataPath string) (err error) {
	store.filename = path.Join(dataPath, store.file)
	if err = util.DecodeJsonFromFile(&kv, store.filename); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	return
}

func (store *kvStore) Save() error {
	return util.EncodeJsonToFile(kv, store.filename)
}

type storeName string

const (
	snNode       storeName = "node"
	snUserMgr    storeName = "user_mgr"
	snCmdMgr     storeName = "cmd_mgr"
	snProcessMgr storeName = "process_mgr"
	snKV         storeName = "kv"
)

var (
	userMgr    *UserMgr
	nodes      = map[string]*Node{}
	cmdMgr     *CmdMgr
	processMgr *ProcessMgr
	kv         = map[string]string{}
)

func init() {
	stores[snNode] = &nodeStore{storeBase{
		file: "node.json",
	}}
	stores[snUserMgr] = &userMgrStore{storeBase{
		file: "user_mgr.json",
	}}
	stores[snCmdMgr] = &cmdMgrStore{storeBase{
		file: "cmd_mgr.json",
	}}
	stores[snProcessMgr] = &processMgrStore{storeBase{
		file: "process_mgr.json",
	}}
	stores[snKV] = &kvStore{storeBase{
		file: "kv.json",
	}}
}
