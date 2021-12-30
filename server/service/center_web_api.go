package service

import (
	"time"
)

func GetNode(cb func(nodes []Node)) {
	centerTaskQueue.Submit(func() {
		nodes := make([]Node, 0)
		time.Sleep(time.Second * 10)
		webTransQueue.Submit(cb, nodes)
	})
}
