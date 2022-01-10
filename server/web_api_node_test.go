package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestNodeHandler_List(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/node/list", address), "POST")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())
	req2, _ = req2.WriteJSON(struct {
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
	}{PageNo: 1, PageSize: 10})

	ret, err := req2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
