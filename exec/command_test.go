package exec

import (
	"bytes"
	"math/rand"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestCommand_Run(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(5) // + 5
	t.Log(num)

	pros := map[int]*Cmd{}

	for i := 0; i < num; i++ {
		//file, err := exec.LookPath("go")
		//t.Log(file, err)
		cmd := Command("./test/test")
		if err := cmd.Run(0, func(cmd *Cmd, err error) {
			// 异步调用过来的
			if err != nil {
				// exit or signal
				t.Error(cmd.Pid(), err)
			} else {
				t.Log(cmd.Pid(), "success")
			}
		}); err != nil {
			t.Error(err)
		} else {
			pros[i] = cmd
		}
	}

	time.Sleep(time.Second)
}

func TestCommandWithCmd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(5) // + 5
	t.Log(num)

	pros := map[int]*Cmd{}
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		ecmd := exec.Command("./test/test")
		errBuff := bytes.Buffer{}
		ecmd.Stderr = &errBuff
		outBuff := bytes.Buffer{}
		ecmd.Stdout = &outBuff
		cmd := CommandWithCmd(ecmd)
		pros[i] = cmd
		wg.Add(1)
		if err := cmd.Run(0, func(cmd *Cmd, err error) {
			wg.Done()
			// 异步调用过来的
			if err != nil {
				// exit or signal
				if cmd.ProcessState().Exited() {
					t.Error("exited", err, errBuff.String())
				} else {
					t.Error("signal", err, errBuff.String())
				}
			} else {
				t.Log("success", errBuff.String(), outBuff.String())
			}
		}); err != nil {
			wg.Done()
			t.Error(err)
		}
	}

	for {
		time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
		find := false
		for _, p := range pros {
			if !p.Done() {
				t.Log(p.Pid(), "kill")
				if err := p.Signal(syscall.SIGTERM); err != nil {
					t.Log(p.Pid(), "kill", err)
				}
				find = true
				break
			}
		}
		if !find {
			t.Log("all done")
			break
		}
	}

	wg.Wait()

}

func TestShell(t *testing.T) {
	shell := "set -euv;mkdir sh;sleep 3s;echo ok"
	ecmd := exec.Command("/bin/sh", "-c", shell)
	ecmd.Dir = "./test"
	outBuff := bytes.Buffer{}
	//errBuff := bytes.Buffer{}
	ecmd.Stderr = &outBuff
	ecmd.Stdout = &outBuff
	cmd := CommandWithCmd(ecmd)
	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := cmd.Run(2, func(cmd *Cmd, err error) {
		// 异步调用过来的
		if err != nil {
			t.Error(err, "--", outBuff.String())
			// exit or signal
			if cmd.ProcessState().Exited() {
				t.Error(outBuff.String())
			} else {
				t.Error(err)
			}
		} else {
			t.Log("success", outBuff.String())
		}
		wg.Done()
	}); err != nil {
		t.Error(err)
		wg.Done()
	}

	wg.Wait()
}
