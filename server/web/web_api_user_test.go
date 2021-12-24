package web

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

func TestUser_Nav(t *testing.T) {
	go startWebListener(t)
	time.Sleep(time.Millisecond * 100)

	ret := authLogin(t, "test", "test")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	req2, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/user/nav", address), "POST")
	req2.SetHeader("token", gjson.Get(ret, "data.token").String())

	ret, err := req2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
