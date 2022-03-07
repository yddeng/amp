package server

import (
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/server/flypd/console/http"
	sproto "github.com/sniperHW/flyfish/server/proto"
	"log"
)

type flyfishHandler struct {
}

var flyClients = map[string]*http.Client{}

func flyCall(host string, req, resp proto.Message, callback func(resp proto.Message, err error)) {
	c, ok := flyClients[host]
	if !ok {
		c = http.NewClient(host)
		flyClients[host] = c
	}
	go func() {
		if _, err := c.Call(req, resp); err != nil {
			taskQueue.Submit(func() {
				callback(nil, err)
			})
		} else {
			taskQueue.Submit(func() {
				callback(resp, nil)
			})
		}
	}()
}

func (*flyfishHandler) GetMeta(done *Done, user string, req struct {
	Host string `json:"host"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, &sproto.GetMeta{}, &sproto.GetMetaResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.GetMetaResp)
		done.result.Data = struct {
			Version int64  `json:"version"`
			Meta    string `json:"meta"`
		}{Version: resp.GetVersion(), Meta: string(resp.GetMeta())}
	})
}

func (*flyfishHandler) AddTable(done *Done, user string, req struct {
	Host    string `json:"host"`
	Name    string `json:"name"`
	Version int64  `json:"version"`
	Fields  []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Default string `json:"default"`
	} `json:"fields"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	msg := &sproto.MetaAddTable{
		Name:    req.Name,
		Version: req.Version,
	}
	for _, f := range req.Fields {
		msg.Fields = append(msg.Fields, &sproto.MetaFiled{
			Name:    f.Name,
			Type:    f.Type,
			Default: f.Default,
		})
	}

	flyCall(req.Host, msg, &sproto.MetaAddTableResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.MetaAddTableResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})

}

func (*flyfishHandler) AddField(done *Done, user string, req struct {
	Host    string `json:"host"`
	Name    string `json:"name"`
	Version int64  `json:"version"`
	Fields  []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Default string `json:"default"`
	} `json:"fields"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	msg := &sproto.MetaAddFields{
		Table:   req.Name,
		Version: req.Version,
	}
	for _, f := range req.Fields {
		msg.Fields = append(msg.Fields, &sproto.MetaFiled{
			Name:    f.Name,
			Type:    f.Type,
			Default: f.Default,
		})
	}

	flyCall(req.Host, msg, &sproto.MetaAddFieldsResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.MetaAddFieldsResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) GetSetStatus(done *Done, user string, req struct {
	Host string `json:"host"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, &sproto.GetSetStatus{}, &sproto.GetSetStatusResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.GetSetStatusResp)
		done.result.Data = resp
	})
}

func (*flyfishHandler) AddSet(done *Done, user string, req struct {
	Host   string         `json:"host"`
	AddSet *sproto.AddSet `json:"addSet"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.AddSet, &sproto.AddSetResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.AddSetResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) RemSet(done *Done, user string, req struct {
	Host   string         `json:"host"`
	RemSet *sproto.RemSet `json:"remSet"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.RemSet, &sproto.RemSetResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.RemSetResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) AddNode(done *Done, user string, req struct {
	Host    string          `json:"host"`
	AddNode *sproto.AddNode `json:"addNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.AddNode, &sproto.AddNodeResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.AddNodeResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) RemNode(done *Done, user string, req struct {
	Host    string          `json:"host"`
	RemNode *sproto.RemNode `json:"remNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.RemNode, &sproto.RemNodeResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.RemNodeResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) AddLeaderStoreToNode(done *Done, user string, req struct {
	Host                  string                        `json:"host"`
	AddLearnerStoreToNode *sproto.AddLearnerStoreToNode `json:"addLearnerStoreToNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.AddLearnerStoreToNode, &sproto.AddLearnerStoreToNodeResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.AddLearnerStoreToNodeResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) RemoveNodeStore(done *Done, user string, req struct {
	Host            string                  `json:"host"`
	RemoveNodeStore *sproto.RemoveNodeStore `json:"removeNodeStore"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.RemoveNodeStore, &sproto.RemoveNodeStoreResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.RemoveNodeStoreResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}

func (*flyfishHandler) SetMarkClear(done *Done, user string, req struct {
	Host         string               `json:"host"`
	SetMarkClear *sproto.SetMarkClear `json:"setMarkClear"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	flyCall(req.Host, req.SetMarkClear, &sproto.SetMarkClearResp{}, func(msg proto.Message, err error) {
		defer func() { done.Done() }()
		if err != nil {
			done.result.Message = err.Error()
			return
		}
		resp := msg.(*sproto.SetMarkClearResp)
		if !resp.GetOk() {
			done.result.Message = resp.GetReason()
		}
	})
}
