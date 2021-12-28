package main

import (
	"flag"
	"initial-sever/logger"
	"initial-sever/server/center"
	"initial-sever/server/web"
	"initial-sever/util"
)

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

var (
	file = flag.String("file", "./config.json", "config")
)

func main() {
	flag.Parse()

	log := logger.NewZapLogger("server.log", "log", "debug", 100, 14, 1, true)
	logger.InitLogger(log)

	var err error
	var cfg Config
	if err = util.DecodeJsonFromFile(&cfg, *file); err != nil {
		panic(err)
	}

	if err = center.LoadData(cfg.DataPath); err != nil {
		panic(err)
	}
	center.RunCenter(cfg.CenterAddress)

	// web
	if err = web.LoadNav(cfg.NavPath); err != nil {
		panic(err)
	}

	if err = web.LoadData(cfg.DataPath, struct {
		Username string
		Password string
	}{Username: cfg.Admin.Username, Password: cfg.Admin.Password}); err != nil {
		panic(err)
	}

	web.RunWeb(cfg.WebAddress)
}
