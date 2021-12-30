package service

import (
	"errors"
	"github.com/yddeng/dnet/drpc"
	"initial-sever/logger"
	"initial-sever/protocol"
)

func (this *Center) onLogin(replier *drpc.Replier, req interface{}) {
	channel := replier.Channel
	msg := req.(*protocol.LoginReq)
	logger.GetSugar().Infof("onLogin %v\n", msg)

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

	client.session = channel.(*Node).session
	client.session.SetContext(client)
	logger.GetSugar().Infof("onLogin %s\n", client.session.RemoteAddr().String())
	replier.Reply(&protocol.LoginResp{}, nil)
}
