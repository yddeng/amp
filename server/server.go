package main

import (
	"encoding/json"
	"flag"
	"initial-sever/server/web"
	"initial-sever/util"
	"io/ioutil"
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

func load(filename string) Config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
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
