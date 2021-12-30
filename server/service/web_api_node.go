package service

import (
	"log"
)

type nodeHandler struct{}

func (*nodeHandler) List(done *Done, user string) {
	log.Printf("info user:%s \n", user)

	GetNode(func(nodes []Node) {
		done.Done()
	})
}
