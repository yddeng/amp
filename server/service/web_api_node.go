package service

import (
	"log"
)

type nodeHandler struct{}

func (*nodeHandler) List(done *Done, user string) {
	log.Printf("info user:%s \n", user)

	getNodeInfo(func(nodes []nodeInfo) {
		done.result.Data = nodes
		done.Done()
	})
}
