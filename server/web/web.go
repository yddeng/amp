package web

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"log"
)

type Config struct {
	WebAddress string `json:"web_address"`
	FilePath   string `json:"file_path"`
}

func RunWeb(address string) {
	/*
	 所有的公共变量在队列中执行。
	 使用warp函数处理过的方法，已经是在队列中执行。
	 方法中不能再调用队列，否则将造成死锁。
	*/

	app := iris.New()
	app.Use(logger.New())
	// 跨域
	app.Use(handleCORS)

	authHandle := new(Auth)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", warpJsonHandle(authHandle.Login))
	authRouter.Post("/logout", warpTokenHandle(authHandle.Logout))

	userHandle := new(User)
	userRouter := app.Party("/user")
	userRouter.Post("/nav", warpTokenHandle(userHandle.Nav))
	userRouter.Post("/info", warpTokenHandle(userHandle.Info))

	log.Printf("web server run %s.\n", address)
	if err := app.Listen(address); err != nil {
		log.Println(err)
	}
}
