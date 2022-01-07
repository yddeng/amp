package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestUserHandler_Nav(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/nav", address), "GET")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())

	ret, err := req2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}

func TestUserHandler_Info(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/info", address), "GET")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())

	ret, err := req2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}

func TestUserHandler_List(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/list", address), "POST")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())
	req2, _ = req2.WriteJSON(struct {
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
	}{PageNo: 1, PageSize: 10})

	ret, _ = req2.ToString()
	t.Log(ret)
}

func TestUserHandler_Add(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/add", address), "POST")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())
	req2, _ = req2.WriteJSON(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: "test", Password: "test"})

	ret, _ = req2.ToString()
	t.Log(ret)
}

func TestUserHandler_Delete(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/delete", address), "POST")
	req2.SetHeader("Access-Token", gjson.Get(ret, "data.token").String())
	req2, _ = req2.WriteJSON(struct {
		Username []string `json:"username"`
	}{Username: []string{"test"}})

	ret, _ = req2.ToString()
	t.Log(ret)
}
