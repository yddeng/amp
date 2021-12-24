package web

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

func TestHandler_GetNav(t *testing.T) {
	go startWebListener()
	time.Sleep(time.Millisecond * 100)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/getNav", address), "POST")
	req2.SetHeader("token", gjson.Get(ret, "data.token").String())

	ret, err := req2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
