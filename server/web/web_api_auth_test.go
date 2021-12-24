package web

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

var address = "127.0.0.1:40235"

func TestLogin(t *testing.T) {
	go func() {
		RunWeb(&Config{
			WebAddress: address,
			FilePath:   "",
		})
	}()

	time.Sleep(time.Second)

	baseUrl := fmt.Sprintf("http://%s", address)
	req, _ := dhttp.PostJson(baseUrl+"/auth/login", struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: "admin", Password: "123456"})

	var ret Result
	if err := req.ToJSON(&ret); err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
