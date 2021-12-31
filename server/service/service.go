package service

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"log"
	"os"
	"time"
)

func NowUnix() int64 {
	return time.Now().Unix()
}

func webRun(address string) {
	/*
	 所有的公共变量在队列中执行。
	 使用warp函数处理过的方法，已经是在队列中执行。
	*/

	app := iris.New()
	app.Use(logger.New())
	// 跨域
	app.Use(handleCORS)

	app.Get("/test", func(ctx iris.Context) {
		var ret Result
		ret.Data = struct {
			Text string `json:"text"`
		}{Text: "hello world!"}
		_, _ = ctx.JSON(ret)
	})

	initHandler(app)

	log.Printf("web server run %s.\n", address)
	go func() {
		if err := app.Listen(address); err != nil {
			panic(err)
		}
	}()
}

func centerRun(address string) {
	center = newCenter(address)
	center.start()
}

type Config struct {
	WebAddress    string `json:"web_address"`
	CenterAddress string `json:"center_address"`
	DataPath      string `json:"data_path"`
	NavPath       string `json:"nav_path"`
	Admin         struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"admin"`
}

func Service(cfg Config) (err error) {
	if err = loadNav(cfg.NavPath); err != nil {
		return
	}

	_ = os.MkdirAll(cfg.DataPath, os.ModePerm)
	if err = loadStore(cfg.DataPath); err != nil {
		return
	}

	if admin == nil {
		admin = &User{
			Username: cfg.Admin.Username,
			Password: cfg.Admin.Password,
		}
		saveStore(snAdmin)
	}

	centerRun(cfg.CenterAddress)
	webRun(cfg.WebAddress)
	return
}
