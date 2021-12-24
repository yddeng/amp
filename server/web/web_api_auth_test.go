package web

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

var address = "127.0.0.1:40235"

func startWebListener() {
	RunWeb(&Config{
		WebAddress: address,
		FilePath:   "",
	})
}

func authLogin(t *testing.T, Username, Password string) string {
	req, _ := dhttp.PostJson(fmt.Sprintf("http://%s/auth/login", address), struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: Username, Password: Password})

	ret, err := req.ToString()
	if err != nil {
		t.Fatal(err)
	}
	return ret
}

func TestLogin(t *testing.T) {
	go startWebListener()
	time.Sleep(time.Millisecond * 100)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret)
}
