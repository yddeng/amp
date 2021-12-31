package service

import (
	"errors"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"initial-sever/logger"
	"initial-sever/protocol"
)

var (
	nodes = map[string]*Node{}
)

type Node struct {
	Name    string       `json:"name"`
	Inet    string       `json:"inet"`
	Net     string       `json:"net"`
	LoginAt int64        `json:"login_at"` // 登陆时间
	session dnet.Session `json:"_"`
}

func (n *Node) Online() bool {
	return n.session != nil
}

func (n *Node) SendRequest(req *drpc.Request) error {
	return n.session.Send(req)
}

func (n *Node) SendResponse(resp *drpc.Response) error {
	return n.session.Send(resp)
}

func (this *Center) onLogin(replier *drpc.Replier, req interface{}) {
	channel := replier.Channel
	msg := req.(*protocol.LoginReq)
	logger.GetSugar().Infof("onLogin %v", msg)

	name := msg.GetName()
	client := nodes[name]
	if client == nil {
		client = &Node{Name: name}
		nodes[name] = client
	}
	if client.session != nil {
		replier.Reply(&protocol.LoginResp{Code: "client already login. "}, nil)
		channel.(*Node).session.Close(errors.New("client already login. "))
		return
	}

	client.Inet = msg.GetInet()
	client.Net = msg.GetNet()
	client.LoginAt = NowUnix()

	client.session = channel.(*Node).session
	client.session.SetContext(client)
	logger.GetSugar().Infof("onLogin %s", client.session.RemoteAddr().String())
	replier.Reply(&protocol.LoginResp{}, nil)
	saveStore(snNode)
}
