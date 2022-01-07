package main

import (
	"flag"
	"initial-server/logger"
	"initial-server/server"
	"initial-server/util"
	"os"
	"os/signal"
	"syscall"
)

var (
	file = flag.String("file", "./center_config.json", "config file")
)

func main() {
	flag.Parse()

	log := logger.NewZapLogger("server.log", "log", "debug", 100, 14, 1, true)
	logger.InitLogger(log)

	var err error
	var cfg server.Config
	if err = util.DecodeJsonFromFile(&cfg, *file); err != nil {
		panic(err)
	}

	if err = server.Service(cfg); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
	}
}
