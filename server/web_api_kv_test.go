package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestKvHandler_Set(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/kv/set", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: "yddeng", Value: "2022"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/kv/get", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Key string `json:"key"`
		}{Key: "yddeng"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}
	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/kv/get", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Key string `json:"key"`
		}{Key: "test"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}
}
