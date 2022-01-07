package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"strings"
	"testing"
)

func TestCmdHandler_List(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()
	{
		//create
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/cmd/create", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Name    string            `json:"name"`
			Dir     string            `json:"dir"`
			Context string            `json:"context"`
			Args    map[string]string `json:"args"`
		}{Name: "test", Dir: "", Context: "sleep 11s;mkdir {{name}};echo ok", Args: map[string]string{"name": "tttt"}})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		//list
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/cmd/list", address), "GET")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			PageNo   int `json:"pageNo"`
			PageSize int `json:"pageSize"`
		}{PageNo: 1, PageSize: 10})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		// exec
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/cmd/exec", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Name    string            `json:"name"`
			Dir     string            `json:"dir"`
			Args    map[string]string `json:"args"`
			Node    string            `json:"node"`
			Timeout int               `json:"timeout"`
		}{Name: "test", Args: map[string]string{"name": "tttt"}, Node: "executor", Timeout: 12})

		ret, err := req.ToString()
		t.Log(err, ret)
		out := gjson.Get(ret, "data.output").String()
		lines := strings.Split(out, "\n")
		for _, v := range lines {
			fmt.Println(v)
		}
	}
}
