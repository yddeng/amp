package task

import (
	"github.com/yddeng/utils/task"
	"reflect"
	"sync"
)

var (
	taskPool   *task.TaskPool
	createOnce sync.Once
)

func Submit(f func()) error {
	createOnce.Do(func() {
		taskPool = task.NewTaskPool(1, 2048)
	})
	return taskPool.Submit(f)
}

func Wait(fn interface{}, args ...interface{}) (ret []interface{}) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != len(args) {
		panic("func argument error")
	}

	vals := make([]reflect.Value, len(args))
	for i, v := range args {
		vals[i] = reflect.ValueOf(v)
	}

	retCh := make(chan []reflect.Value)
	if err := Submit(func() {
		retCh <- val.Call(vals)
	}); err != nil {
		panic(err)
	}

	vals = <-retCh
	for _, v := range vals {
		ret = append(ret, v.Interface())
	}
	return
}
