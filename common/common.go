package common

// 隐藏目录，存放临时配置文件
const dir = ".amp"

// 程序状态
const (
	StateUnknown  = "Unknown"
	StateStarting = "Starting"
	StateRunning  = "Running"
	StateStopping = "Stopping"
	StateStopped  = "Stopped"
	StateExited   = "Exited"
)
