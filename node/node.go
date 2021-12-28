package main

import (
	"initial-sever/logger"
	"initial-sever/node/client"
)

func main() {
	log := logger.NewZapLogger("client.log", "log", "debug", 100, 14, 1, true)
	logger.InitLogger(log)
	c := client.NewClient(client.Config{
		Name:     "client",
		Center:   "127.0.0.1:40155",
		FilePath: "./data",
	})

	c.Start()

	select {}
}
