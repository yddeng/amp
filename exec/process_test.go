package exec

import (
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestProcessWithPid(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	num := 5 //rand.Intn(5) + 5
	t.Log(num)

	pros := map[int]*Cmd{}
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		ecmd := exec.Command("./test/test", "./test/config.json")
		errBuff := bytes.Buffer{}
		ecmd.Stderr = &errBuff
		outBuff := bytes.Buffer{}
		ecmd.Stdout = &outBuff
		cmd := CommandWithCmd(ecmd)
		pros[i] = cmd
		wg.Add(1)
		if err := cmd.cmd.Start(); err != nil {
			wg.Done()
			t.Error(err)
		} else {
			p, err := os.FindProcess(cmd.Pid())
			if err != nil {
				t.Error(err)
			}
			t.Log("findProcess", p.Pid)
			go func() {
				state, err := p.Wait()
				close(cmd.done)
				wg.Done()
				t.Log(state.Pid(), err, state.ExitCode(), state.Exited(), state.Success(), state.String())
			}()
		}
	}

	for {
		time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
		find := false
		for _, p := range pros {
			if !p.Done() {
				sig := syscall.SIGTERM
				code := rand.Intn(3)
				if code == 1 {
					sig = syscall.SIGTERM
				} else if code == 2 {
					sig = syscall.SIGKILL
				}
				t.Log(p.Pid(), "signal", sig.String())
				if err := p.Signal(sig); err != nil {
					t.Log(p.Pid(), "signal", err)
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
