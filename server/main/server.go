package main

import (
	"initial-sever/server/web"
	"os"
)

func main() {
	address := os.Args[1]
	web.RunWeb(&web.Config{
		WebAddress: address,
	})
}
