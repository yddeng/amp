package web

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"log"
	"time"
)

type token struct {
	user   string
	expire time.Time
}

var (
	tknUser = map[string]*token{}
	dur     = time.Hour
)

func addToken(u string) (tkn string) {
	now := time.Now()
	tkn = fmt.Sprintf("%d", now.UnixNano())
	t := &token{
		user:   u,
		expire: now.Add(dur),
	}
	tknUser[tkn] = t
	return
}

func getToken(tkn string) (string, bool) {
	t, ok := tknUser[tkn]
	if !ok {
		return "", false
	}
	t.expire = time.Now().Add(dur)
	return t.user, true
}

type Config struct {
	WebAddress string `json:"web_address"`
	FilePath   string `json:"file_path"`
}

func RunWeb(address string) {
	app := iris.New()
	app.Use(logger.New())
	// 跨域
	app.Use(handleCORS)

	authHandle := new(Auth)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", warpJsonHandle(authHandle.Login))

	userHandle := new(User)
	userRouter := app.Party("/user")
	userRouter.Post("/nav", warpTokenHandle(userHandle.Nav))

	log.Printf("web server run %s.\n", address)
	if err := app.Listen(address); err != nil {
		panic(err)
	}
}
