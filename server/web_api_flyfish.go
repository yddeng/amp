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
