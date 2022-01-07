package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(4)
	if num == 0 {
		panic("panic")
	} else if num == 1 {
		return
	} else if num == 2 {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 65535)
				l := runtime.Stack(buf, false)
				fmt.Println(fmt.Sprintf("%v: %s", r, buf[:l]))
			}
		}()
		panic("panic and recover")
	} else {
		pid := os.Getpid()
		i := 1
		for {
			time.Sleep(time.Millisecond * 100)
			i++
			fmt.Println(pid, "---", i)
		}
	}
}
