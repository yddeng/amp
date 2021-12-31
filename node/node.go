package main

import (
	"initial-server/logger"
	"initial-server/node/client"
)

func main() {
	log := logger.NewZapLogger("client.log", "log", "debug", 100, 14, 1, true)
	logger.InitLogger(log)
	c := client.NewClient(client.Config{
		Name:     "client",
		Net:      "",
		Inet:     "10.128.2.123",
		Center:   "10.128.2.123:40155",
		FilePath: "./data",
	})

	c.Start()

	select {}
}
