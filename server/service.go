package server

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/yddeng/utils/task"
	"log"
	"os"
	"time"
)

func NowUnix() int64 {
	return time.Now().Unix()
}

type Config struct {
	DataPath     string        `json:"data_path"`
	CenterConfig *CenterConfig `json:"center_config"`
	WebConfig    *WebConfig    `json:"web_config"`
}

type CenterConfig struct {
	Address string `json:"address"`
}

type WebConfig struct {
	Address string `json:"address"`
	Admin   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"admin"`
	Nav []*Nav `json:"nav"`
}

var (
	dataPath  string
	taskQueue = task.NewTaskPool(1, 2048)
)

func Service(cfg Config) (err error) {
	dataPath = cfg.DataPath
	_ = os.MkdirAll(cfg.DataPath, os.ModePerm)
	if err = loadStore(cfg.DataPath); err != nil {
		return
	}

	centerRun(cfg.CenterConfig)
	webRun(cfg.WebConfig)
	return
}

func webRun(cfg *WebConfig) {
	allNav = cfg.Nav
	if admin == nil {
		admin = &User{
			Username: cfg.Admin.Username,
			Password: cfg.Admin.Password,
		}
		saveStore(snAdmin)
	}

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

	log.Printf("web server run %s.\n", cfg.Address)
	go func() {
		if err := app.Listen(cfg.Address); err != nil {
			panic(err)
		}
	}()
}

func centerRun(cfg *CenterConfig) {
	center = newCenter(cfg.Address)
	center.start()
}
