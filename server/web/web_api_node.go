package web

import (
	"initial-sever/server/center"
	"log"
)

type Node struct {
}

func (*Node) List(user string) (ret Result) {
	log.Printf("info user:%s \n", user)
	ret.Data = center.GetNodes()
	return
}
