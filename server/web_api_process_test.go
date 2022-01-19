package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestProcessHandler_Group(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/gadd", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Group string `json:"group"`
		}{Group: "all"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/gadd", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Group string `json:"group"`
		}{Group: "all/test"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/gadd", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Group string `json:"group"`
		}{Group: "nav/nav1/nav2"})
		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/glist", address), "POST")
		req.SetHeader("Access-Token", token)
		ret, err := req.ToString()
		t.Log(err, ret)
	}

	//{
	//	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/gremove", address), "POST")
	//	req.SetHeader("Access-Token", token)
	//	req.WriteJSON(struct {
	//		Group string `json:"group"`
	//	}{Group: "all"})
	//	ret, err := req.ToString()
	//	t.Log(err, ret)
	//}
	//
	//{
	//	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/gremove", address), "POST")
	//	req.SetHeader("Access-Token", token)
	//	req.WriteJSON(struct {
	//		Group string `json:"group"`
	//	}{Group: "nav/nav1"})
	//	ret, err := req.ToString()
	//	t.Log(err, ret)
	//}

	{
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/glist", address), "POST")
		req.SetHeader("Access-Token", token)
		ret, err := req.ToString()
		t.Log(err, ret)
	}
}

func TestProcessHandler_List(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	{
		//create
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/create", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Dir          string           `json:"dir"`
			Config       []*ProcessConfig `json:"config"`
			Command      string           `json:"command"`
			Priority     int              `json:"priority"`
			StartRetries int              `json:"start_retries"`
			StopWaitSecs int              `json:"stop_wait_secs"`
			Groups       []string         `json:"groups"`
			Node         string           `json:"node"`
		}{Dir: "/Users/yidongdeng/go/src/initial-server/exec/test",
			Config: []*ProcessConfig{
				{Name: "config.json", Context: `{"msg":"teststetstest"}`},
			},
			Command:      "./test {{id}}/config.json",
			Priority:     1,
			StartRetries: 1,
			StopWaitSecs: 1,
			Groups:       []string{"hhh"},
			Node:         "executor",
		})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	/*
		{
			// exec
			req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/start", address), "POST")
			req.SetHeader("Access-Token", token)
			req.WriteJSON(struct {
				ID int `json:"id"`
			}{ID: 1})

			ret, err := req.ToString()
			t.Log(err, ret)
			out := gjson.Get(ret, "data.result").String()
			lines := strings.Split(out, "\n")
			for _, v := range lines {
				fmt.Println(v)
			}
		}

	*/

}
