package server

import (
	"errors"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"initial-server/logger"
	"initial-server/protocol"
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
	if n.session == nil {
		return errors.New("session is nil")
	}
	return n.session.Send(req)
}

func (n *Node) SendResponse(resp *drpc.Response) error {
	if n.session == nil {
		return errors.New("session is nil")
	}
	return n.session.Send(resp)
}

func (this *Center) onLogin(replier *drpc.Replier, req interface{}) {
	channel := replier.Channel
	msg := req.(*protocol.LoginReq)
	logger.GetSugar().Infof("onLogin %v", msg)

	if this.token != "" && msg.GetToken() != this.token {
		replier.Reply(&protocol.LoginResp{Code: "token failed"}, nil)
		channel.(*Node).session.Close(errors.New("token failed. "))
		return
	}

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
