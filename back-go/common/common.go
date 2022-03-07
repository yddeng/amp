package common

import "time"

// 隐藏目录，存放临时配置文件
const AmpDir = ".amp"

// 程序状态
const (
	StateUnknown  = "Unknown"
	StateStarting = "Starting"
	StateRunning  = "Running"
	StateStopping = "Stopping"
	StateStopped  = "Stopped"
	StateExited   = "Exited"
)

const HeartbeatTimeout = 10 * time.Second
