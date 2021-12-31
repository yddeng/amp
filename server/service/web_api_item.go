package service

import (
	"log"
	"sort"
)

type itemHandler struct {
}

func (*itemHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("item/list by(%s) %v\n", user, req)

	getItemList(func(items []*Item) {
		sort.Slice(items, func(i, j int) bool {
			return items[i].UpdateAt > items[j].UpdateAt
		})
		start, end := listRange(req.PageNo, req.PageSize, len(items))
		done.result.Data = items[start:end]
		done.Done()
	})
}
