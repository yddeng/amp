package service

import (
	"log"
)

type Node struct {
}

func (*Node) List(done *Done, user string) {
	log.Printf("info user:%s \n", user)

	GetNode(func(nodes []Node) {
		done.Done()
	})
}
