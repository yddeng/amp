package main

import (
	"amp/back-go/server"
	"amp/back-go/util"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	file = flag.String("file", "./amps_config.json", "config file")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	//log := logger.NewZapLogger("server.log", "log", "debug", 100, 14, 1, true)
	//logger.InitLogger(log)

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
	log.Printf("receive signal:%v to stopping. ", <-sigChan)
	server.Stop()
	log.Println("stopped. ")
}
