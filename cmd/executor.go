package main

import (
	"initial-server/exec"
	"initial-server/logger"
	"initial-server/server"
	"initial-server/util"
)

func main() {
	log := logger.NewZapLogger("client.log", "log", "debug", 100, 14, 1, true)
	logger.InitLogger(log)

	var err error
	var cfg exec.Config
	if err = util.DecodeJsonFromFile(&cfg, *file); err != nil {
		panic(err)
	}

	if err = exec.Service(cfg); err != nil {
		panic(err)
	}

	c := exec.NewClient(exec.Config{
		Name:     "client",
		Net:      "",
		Inet:     "10.128.2.123",
		Center:   "10.128.2.123:40155",
		FilePath: "./data",
	})

	c.Start()

	select {}
}
