package main

import (
	"flag"
	"initial-sever/server/web"
	"initial-sever/util"
)

type Config struct {
	WebAddress string `json:"web_address"`
	DataPath   string `json:"data_path"`
	NavPath    string `json:"nav_path"`
	Admin      struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"admin"`
}

var (
	file = flag.String("file", "./config.json", "config")
)

func main() {
	flag.Parse()

	var cfg Config
	if err := util.DecodeJsonFromFile(&cfg, *file); err != nil {
		panic(err)
	}

	if err := web.LoadNav(cfg.NavPath); err != nil {
		panic(err)
	}

	web.RunWeb(cfg.WebAddress)
}
