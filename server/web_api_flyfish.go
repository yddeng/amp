package server

import (
	"github.com/sniperHW/flyfish/server/flypd/console/http"
	sproto "github.com/sniperHW/flyfish/server/proto"
	"log"
)

type flyfishHandler struct {
}

var flyClients = map[string]*http.Client{}

func getFlyClient(host string) *http.Client {
	c, ok := flyClients[host]
	if !ok {
		c = http.NewClient(host)
		flyClients[host] = c
	}
	return c
}

func (*flyfishHandler) GetMeta(done *Done, user string, req struct {
	Host string `json:"host"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.GetMetaResp
	if _, err := c.Call(&sproto.GetMeta{}, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	done.result.Data = struct {
		Version int64  `json:"version"`
		Meta    string `json:"meta"`
	}{Version: resp.GetVersion(), Meta: string(resp.GetMeta())}
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
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

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

	var resp sproto.MetaAddTableResp
	if _, err := c.Call(msg, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
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
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

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

	var resp sproto.MetaAddTableResp
	if _, err := c.Call(msg, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) GetSetStatus(done *Done, user string, req struct {
	Host string `json:"host"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.GetSetStatusResp
	if _, err := c.Call(&sproto.GetSetStatus{}, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	done.result.Data = resp
}

func (*flyfishHandler) AddSet(done *Done, user string, req struct {
	Host   string         `json:"host"`
	AddSet *sproto.AddSet `json:"addSet"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.AddSetResp
	if _, err := c.Call(req.AddSet, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) RemSet(done *Done, user string, req struct {
	Host   string         `json:"host"`
	RemSet *sproto.RemSet `json:"remSet"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.RemSetResp
	if _, err := c.Call(req.RemSet, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) AddNode(done *Done, user string, req struct {
	Host    string          `json:"host"`
	AddNode *sproto.AddNode `json:"addNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.AddNodeResp
	if _, err := c.Call(req.AddNode, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) RemNode(done *Done, user string, req struct {
	Host    string          `json:"host"`
	RemNode *sproto.RemNode `json:"remNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.RemNodeResp
	if _, err := c.Call(req.RemNode, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) AddLeaderStoreToNode(done *Done, user string, req struct {
	Host                  string                        `json:"host"`
	AddLearnerStoreToNode *sproto.AddLearnerStoreToNode `json:"addLearnerStoreToNode"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.AddLearnerStoreToNodeResp
	if _, err := c.Call(req.AddLearnerStoreToNode, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) RemoveNodeStore(done *Done, user string, req struct {
	Host            string                  `json:"host"`
	RemoveNodeStore *sproto.RemoveNodeStore `json:"removeNodeStore"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.RemoveNodeStoreResp
	if _, err := c.Call(req.RemoveNodeStore, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}

func (*flyfishHandler) SetMarkClear(done *Done, user string, req struct {
	Host         string               `json:"host"`
	SetMarkClear *sproto.SetMarkClear `json:"setMarkClear"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	c := getFlyClient(req.Host)

	var resp sproto.SetMarkClearResp
	if _, err := c.Call(req.SetMarkClear, &resp); err != nil {
		done.result.Message = err.Error()
		return
	}

	if !resp.GetOk() {
		done.result.Message = resp.GetReason()
	}
}
