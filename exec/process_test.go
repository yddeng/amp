package exec

import (
	"os/exec"
	"sync"
	"testing"
)

func TestNewProcess(t *testing.T) {

	cmd := exec.Command("./test/test")
	if err := cmd.Start(); err != nil {
		t.Error(err)
		return
	}

	wg := sync.WaitGroup{}

	t.Log("process start", cmd.Process.Pid)
	p, err := NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		t.Error(err)
		return
	}

	wg.Add(1)
	p.waitCmd(cmd, func(process *Process) {
		t.Log(process.Pid, process.GetState())
		wg.Done()
	})

	wg.Wait()
}

func TestProcessWait(t *testing.T) {

	p, err := NewProcess(96543)
	if err != nil {
		t.Error(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	p.waitNoChild(func(process *Process) {
		t.Log(process.Pid, process.GetState())
		wg.Done()
	})

	wg.Wait()
}
