package server

import (
	"log"
	"sort"
)

type nodeHandler struct{}

func (*nodeHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	getNodeInfo(func(nodes []*nodeInfo) {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].LoginAt > nodes[j].LoginAt
		})
		start, end := listRange(req.PageNo, req.PageSize, len(nodes))
		done.result.Data = struct {
			PageNo     int         `json:"pageNo"`
			PageSize   int         `json:"pageSize"`
			TotalCount int         `json:"totalCount"`
			NodeList   []*nodeInfo `json:"nodeList"`
		}{PageNo: req.PageNo,
			PageSize:   req.PageSize,
			TotalCount: len(nodes),
			NodeList:   nodes[start:end],
		}
		done.Done()
	})
}
