package exec

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
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
		ecmd := exec.Command("./test/test", "./test/config.json")
		errBuff := bytes.Buffer{}
		ecmd.Stderr = &errBuff
		outBuff := bytes.Buffer{}
		ecmd.Stdout = &outBuff
		cmd := CommandWithCmd(ecmd)
		pros[i] = cmd
		wg.Add(1)
		if err := cmd.Run(0, func(cmd *Cmd, err error) {
			wg.Done()
			t.Log(err, cmd.cmd.ProcessState.ExitCode(), cmd.ProcessState().Exited(), cmd.ProcessState().Success())
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

func TestCmd_Pid(t *testing.T) {
	ecmd := exec.Command("/bin/sh", "-c", "./test")
	ecmd.Dir = "./test"

	f, err := os.OpenFile("test_err.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
	//errout := bytes.Buffer{}
	ecmd.Stderr = f
	ecmd.Stdout = f
	defer f.Close()

	err = ecmd.Start()
	if err != nil {
		t.Error(err)
	}

	t.Log(ecmd.Process.Pid, ecmd.ProcessState)
	time.AfterFunc(time.Second, func() {
		syscall.Kill(0, syscall.SIGTERM)
		os.Exit(0)
	})
	err = ecmd.Wait()

	data, err := ioutil.ReadFile("test_err.log")
	if err != nil {
		t.Error(err)
	}

	t.Log("end", err, "-", string(data), "-")

}
