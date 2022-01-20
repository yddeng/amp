package main

import (
	"initial-server/logger"
	"initial-server/util"
	"os"
)

type Config struct {
	Msg string `json:"msg"`
}

func main() {

	filename := os.Args[1]
	var cfg Config
	if err := util.DecodeJsonFromFile(&cfg, filename); err != nil {
		panic(err)
	}

	logg := logger.NewZapLogger("test.log", "log", "debug", 100, 14, 1, true)
	logg.Sugar().Info(cfg.Msg)

	//rand.Seed(time.Now().UnixNano())
	//num := rand.Intn(4)
	//if num == 0 {
	panic("panic")
	//} else if num == 1 {
	//
	//}
	//	return
	//} else if num == 2 {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			buf := make([]byte, 65535)
	//			l := runtime.Stack(buf, false)
	//			logg.Sugar().Info(fmt.Sprintf("%v: %s", r, buf[:l]))
	//		}
	//	}()
	//	panic("panic and recover")
	//} else {
	//	pid := os.Getpid()
	//	i := 1
	//	for {
	//		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))
	//		i++
	//		logg.Sugar().Info(pid, "---", i)
	//		if i == 20 {
	//			break
	//		}
	//	}
	//}

	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//select {
	//case <-sigChan:
	//}
	//logg.Sugar().Info("listen stopping. ")
	//time.Sleep(time.Second * 10)
	//logg.Sugar().Info("listen stopped ")
}
