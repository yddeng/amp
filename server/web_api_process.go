package server

import (
	"log"
)

type ProcessConfig struct {
	Name    string `json:"name"`
	Context string `json:"context"`
}

const (
	processStarting = "Starting"
	processRunning  = "Running"
	processFatal    = "Fatal"  // 设置重启时，超过重启次数
	processExited   = "Exited" // 程序退出
	processStopping = "Stopping"
	processStopped  = "Stopped"
)

type ProcessState struct {
	Pid     int    `json:"pid"`
	StartAt int64  `json:"start_at"`
	Status  string `json:"status"`
}

type Process struct {
	ID           int              `json:"id"`
	Dir          string           `json:"dir"`
	Config       []*ProcessConfig `json:"config"`
	Command      string           `json:"command"`
	Priority     int              `json:"priority"`       // 子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭	默认值为999 。。非必须设置
	StartRetries int              `json:"start_retries"`  // 当进程启动失败后，最大尝试启动的次数。。当超过3次后，supervisor将把 此进程的状态置为FAIL	默认值为3 。。非必须设置
	StopWaitSecs int              `json:"stop_wait_secs"` // 这个是当我们向子进程发送stopsignal信号后，到系统返回信息	给supervisord，所等待的最大时间。 超过这个时间，supervisord会向该	子进程发送一个强制kill的信号。
	State        ProcessState     `json:"state"`
	User         string           `json:"user"`
	CreateAt     int64            `json:"create_at"`
}

type ProcessGroup struct {
	Process map[int]*Process
}

type ProcessMgr struct {
	GenID   int                      `json:"gen_id"`
	Process map[int]*Process         `json:"process"`
	Group   []string                 `json:"group"` // 程序组
	Groups  map[string]*ProcessGroup `json:"_"`     // 运行时设置
}

type processHandler struct {
}

func (*processHandler) List(done *Done, user string, req struct {
	Group    string `json:"group"`
	PageNo   int    `json:"pageNo"`
	PageSize int    `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

}

func (*processHandler) Create(done *Done, user string, req struct {
	ID           int              `json:"id"`
	Dir          string           `json:"dir"`
	Config       []*ProcessConfig `json:"config"`
	Command      string           `json:"command"`
	Priority     int              `json:"priority"`
	StartRetries int              `json:"start_retries"`
	StopWaitSecs int              `json:"stop_wait_secs"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

}

func (*processHandler) Delete(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

}

func (*processHandler) Group(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

}
