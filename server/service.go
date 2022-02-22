package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/yddeng/utils/task"
	"log"
	"os"
	"strings"
	"time"
)

func NowUnix() int64 {
	return time.Now().Unix()
}

type Config struct {
	DataPath       string        `json:"data_path"`
	CmdLogCapacity int           `json:"cmd_log_capacity"`
	CenterConfig   *CenterConfig `json:"center_config"`
	WebConfig      *WebConfig    `json:"web_config"`
}

type CenterConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

type WebConfig struct {
	Address string `json:"address"`
	App     string `json:"app"`
	Admin   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"admin"`
	Nav []*Nav `json:"nav"`
}

var (
	dataPath  string
	taskQueue = task.NewTaskPool(1, 2048)
	app       *iris.Application
)

func Service(cfg Config) (err error) {
	dataPath = cfg.DataPath
	_ = os.MkdirAll(cfg.DataPath, os.ModePerm)
	if err = loadStore(cfg.DataPath); err != nil {
		return
	}
	cmdLogCapacity = cfg.CmdLogCapacity

	centerRun(cfg.CenterConfig)
	webRun(cfg.WebConfig)
	return
}

func Stop() {
	ch := make(chan bool)
	taskQueue.Submit(func() {
		app.Shutdown(context.TODO())
		saveStore()
		ch <- true
	})
	<-ch
}

func webRun(cfg *WebConfig) {
	allNav = cfg.Nav
	if userMgr.Admin == nil {
		userMgr.Admin = &User{
			Username: cfg.Admin.Username,
			Password: cfg.Admin.Password,
		}
		saveStore(snUserMgr)
	}

	/*
	 所有的公共变量在队列中执行。
	 使用warp函数处理过的方法，已经是在队列中执行。
	*/

	app = iris.New()
	app.Logger().SetLevel("disable")
	// 跨域
	app.Use(handleCORS)

	if cfg.App != "" {
		dir := app.Party("/")
		dir.HandleDir("", cfg.App, iris.DirOptions{
			IndexName: "index.html",
			Gzip:      false,
			ShowList:  false,
		})
	}

	// 代理
	redirect := func(ctx iris.Context) {
		if strings.HasPrefix(ctx.Path(), "/api") {
			name := strings.TrimPrefix(ctx.Path(), "/api")
			if name != "" {
				ctx.Exec(ctx.Method(), name)
				return
			}
		}
	}
	api := app.Party("/api")
	api.Get("/*", redirect)
	api.Post("/*", redirect)

	initHandler(app)

	log.Printf("web server run %s.\n", cfg.Address)
	go func() {
		if err := app.Listen(cfg.Address); err != nil {
			if err == iris.ErrServerClosed {
				log.Printf("web server %s stoped.\n", cfg.Address)
			} else {
				panic(err)
			}
		}
	}()
}

func centerRun(cfg *CenterConfig) {
	log.Printf("center server run %s.\n", cfg.Address)
	center = newCenter(cfg.Address, cfg.Token)
	go func() {
		if err := center.startListener(); err != nil {
			panic(err)
		}
	}()

	go func() {
		timer := time.NewTimer(time.Second)
		for {
			<-timer.C
			taskQueue.Submit(func() {
				processTick()
				processAutoStart()
				timer.Reset(time.Second)
			})
		}
	}()

}
