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

	app.Get("/test", func(context iris.Context) {
		_, _ = context.JSON(Result{Data: struct {
			Text string
		}{Text: "hello world!"}})
	})

	authHandle := new(Auth)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", warpJsonHandle(authHandle.Login))
	authRouter.Post("/logout", warpTokenHandle(authHandle.Logout))

	userHandle := new(User)
	userRouter := app.Party("/user")
	userRouter.Get("/nav", warpTokenHandle(userHandle.Nav))
	userRouter.Get("/info", warpTokenHandle(userHandle.Info))
	userRouter.Post("/list", warpTokenJsonHandle(userHandle.List))
	userRouter.Post("/add", warpTokenJsonHandle(userHandle.Add))
	userRouter.Post("/delete", warpTokenJsonHandle(userHandle.Delete))

	nodeHandle := new(Node)
	nodeRouter := app.Party("/node")
	nodeRouter.Get("/list", warpTokenHandle(nodeHandle.List))

	log.Printf("web server run %s.\n", address)
	if err := app.Listen(address); err != nil {
		log.Println(err)
	}
}
