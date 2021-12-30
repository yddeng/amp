package service

import "github.com/yddeng/utils/task"

var (
	centerTaskQueue = task.NewTaskPool(1, 2048)
)
