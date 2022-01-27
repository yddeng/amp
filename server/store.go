package server

import (
	"amp/util"
	"errors"
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
	for name, store := range stores {
		if err = store.Load(dataPath); err != nil {
			return errors.New(string(name) + " : " + err.Error())
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

func (store *processMgrStore) Save() {
	_ = util.EncodeJsonToFile(processMgr, store.filename)
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
	snNode       storeName = "node"
	snUserMgr    storeName = "user_mgr"
	snCmdMgr     storeName = "cmd_mgr"
	snProcessMgr storeName = "process_mgr"
)

var (
	userMgr    *UserMgr
	nodes      = map[string]*Node{}
	cmdMgr     *CmdMgr
	processMgr *ProcessMgr
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
}
