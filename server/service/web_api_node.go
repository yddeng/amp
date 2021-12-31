package service

import (
	"log"
	"sort"
)

type nodeHandler struct{}

func listRange(pageNo, pageSize, length int) (start int, end int) {
	start = (pageNo - 1) * pageSize
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}
	end = start + pageSize
	if end > length {
		end = length
	}
	return
}

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
		done.result.Data = nodes[start:end]
		done.Done()
	})
}
