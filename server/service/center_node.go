package service

import (
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
)

var (
	nodes map[string]*Node
)

type Node struct {
	Name    string
	session dnet.Session
}

func (c *Node) SendRequest(req *drpc.Request) error {
	return c.session.Send(req)
}

func (c *Node) SendResponse(resp *drpc.Response) error {
	return c.session.Send(resp)
}
