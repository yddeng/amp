package server

import (
	"log"
)

type kvHandler struct {
}

func (*kvHandler) Set(done *Done, user string, req struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	kv[req.Key] = req.Value
	saveStore(snKV)
}

func (*kvHandler) Get(done *Done, user string, req struct {
	Key string `json:"key"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	v, ok := kv[req.Key]
	done.result.Data = struct {
		Value string `json:"value"`
		Exist bool   `json:"exist"`
	}{Value: v, Exist: ok}
}

func (*kvHandler) Delete(done *Done, user string, req struct {
	Key string `json:"key"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if v, ok := kv[req.Key]; ok {
		done.result.Data = struct {
			Value string `json:"value"`
		}{Value: v}
		delete(kv, req.Key)
		saveStore(snKV)
	}
}
