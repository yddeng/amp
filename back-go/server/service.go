package server

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	Token   string `json:"token"`
}

type WebConfig struct {
	Address string `json:"address"`
	App     string `json:"app"`
	Admin   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"admin"`
}

var (
	dataPath  string
	taskQueue = task.NewTaskPool(1, 2048)
	app       *gin.Engine
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

func Stop() {
	ch := make(chan bool)
	taskQueue.Submit(func() {
		doSave(true)
		ch <- true
	})
	<-ch
}

func webRun(cfg *WebConfig) {
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

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	// 前端
	if cfg.App != "" {
		app.Use(static.Serve("/", static.LocalFile(cfg.App, false)))
		app.NoRoute(func(ctx *gin.Context) {
			ctx.File(cfg.App + "/index.html")
		})
	}

	initHandler(app)

	// vue项目路由 /api
	for _, r := range app.Routes() {
		app.Handle(r.Method, "/api"+r.Path, r.HandlerFunc)
	}

	log.Printf("web server run %s.\n", cfg.Address)
	go func() {
		if err := app.Run(cfg.Address); err != nil {
			panic(err)
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
				doSave(false)
				timer.Reset(time.Second)
			})
		}
	}()

}
