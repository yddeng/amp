package web

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

var address = ":40156"

func startWebListener(t *testing.T) {
	if err := LoadNav("../nav.json"); err != nil {
		t.Fatal(err)
	}
	RunWeb(address)
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

func TestAuth_Login(t *testing.T) {
	go startWebListener(t)
	time.Sleep(time.Millisecond * 100)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret)
}
