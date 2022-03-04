package main

import (
	"amp/exec"
	"amp/util"
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	file = flag.String("file", "./ampe_config.json", "config file")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	//log := logger.NewZapLogger("executor.log", "log", "debug", 100, 14, 1, true)
	//logger.InitLogger(log)
	log.Println("ampe")

	var err error
	var cfg exec.Config
	if err = util.DecodeJsonFromFile(&cfg, *file); err != nil {
		panic(err)
	}

	if err = exec.Start(cfg); err != nil {
		panic(err)
	}

	go func() {
		http.ListenAndServe("10.128.2.123:40160", nil)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("receive signal:%v to stopping. ", <-sigChan)
	exec.Stop()
	log.Println("stopped. ")
}
