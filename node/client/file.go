package client

import (
	"initialtool/deploy/logger"
	"initialtool/deploy/util"
	"syscall"
	"time"
)

const (
	dataFile = "execInfo.json"
)

func loadExecInfo(filename string) (execInfos map[int32]*execInfo) {
	delInfos := []int32{}
	if err := util.DecodeJsonFromFile(&execInfos, filename); err == nil {
		for id, info := range execInfos {
			isAlive := info.isAlive()
			logger.GetSugar().Infof("loadExecInfo %v isAlive:%v", info, isAlive)
			if !isAlive {
				delInfos = append(delInfos, id)
			}
		}
	}

	for _, id := range delInfos {
		delete(execInfos, id)
	}
	writeExecInfo(filename, execInfos)
	return
}

func writeExecInfo(filename string, execInfos map[int32]*execInfo) {
	if err := util.EncodeJsonToFile(execInfos, filename); err != nil {
		logger.GetSugar().Errorf(err.Error())
	}
}

type execInfo struct {
	ItemID    int32 `json:"item_id,omitempty"`
	Pid       int   `json:"pid,omitempty"`
	Timestamp int64 `json:"timestamp,omitempty"`
}

func (this *execInfo) isAlive() bool {
	if err := syscall.Kill(this.Pid, 0); err == nil {
		return true
	}
	return false
}

func (c *Client) addExecInfo(itemID int32, pid int) *execInfo {
	info := &execInfo{
		ItemID:    itemID,
		Pid:       pid,
		Timestamp: time.Now().Unix(),
	}

	c.execInfos[itemID] = info
	writeExecInfo(c.execFilename, c.execInfos)
	return info
}

func (c *Client) delExecInfo(execId int32) {
	delete(c.execInfos, execId)
	writeExecInfo(c.execFilename, c.execInfos)
}
