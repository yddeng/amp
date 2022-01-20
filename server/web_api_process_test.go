package server

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
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

func TestProcessHandler_Start(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()
	/*
		{
			//create
			req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/create", address), "POST")
			req.SetHeader("Access-Token", token)
			req.WriteJSON(struct {
				Dir            string           `json:"dir"`
				Config         []*ProcessConfig `json:"config"`
				Command        string           `json:"command"`
				Priority       int              `json:"priority"`
				StopWaitSecs   int              `json:"stop_wait_secs"`
				AutoStartTimes int              `json:"auto_start_times"`
				Groups         []string         `json:"groups"`
				Node           string           `json:"node"`
			}{Dir: "",
				Config: []*ProcessConfig{
					{Name: "config.json", Context: `{"msg":"teststetstest"}`},
				},
				Command:        "./test {{path}}/config.json",
				Priority:       1,
				StopWaitSecs:   5,
				AutoStartTimes: 3,
				Groups:         []string{"hhh"},
				Node:           "executor",
			})

			ret, err := req.ToString()
			t.Log(err, ret)
		}

	*/

	{
		// exec
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/start", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			ID int `json:"id"`
		}{ID: 1})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		running := false
		// list
		for i := 0; i < 3; i++ {
			req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/list", address), "POST")
			req.SetHeader("Access-Token", token)
			req.WriteJSON(struct {
				Group string `json:"group"`
			}{Group: "hhh"})

			ret, err := req.ToString()
			t.Log(err, ret)
			t.Log(gjson.Get(ret, "data.1.state").String())

			running = gjson.Get(ret, "data.1.state.status").String() == "Running"

			time.Sleep(time.Second)
		}

		if running {
			req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/stop", address), "POST")
			req.SetHeader("Access-Token", token)
			req.WriteJSON(struct {
				ID int `json:"id"`
			}{ID: 1})
			ret, err := req.ToString()
			t.Log(err, ret)
		}

		stopping := true
		for i := 0; i < 10 && stopping; i++ {
			req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/process/list", address), "POST")
			req.SetHeader("Access-Token", token)
			req.WriteJSON(struct {
				Group string `json:"group"`
			}{Group: "hhh"})

			ret, err := req.ToString()
			t.Log(err, ret)
			t.Log(gjson.Get(ret, "data.1.state").String())
			stopping = gjson.Get(ret, "data.1.state.status").String() == "Stopping"

			time.Sleep(time.Second)
		}

	}

}
