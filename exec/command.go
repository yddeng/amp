package exec

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Cmd struct {
	cmd    *exec.Cmd
	done   chan struct{}
	doneFn func(cmd *Cmd, err error)
}

func (cmd *Cmd) String() string {
	return cmd.cmd.String()
}

// 异步执行
func (cmd *Cmd) Run(timeout int, fn func(cmd *Cmd, err error)) error {
	cmd.doneFn = fn

	if err := cmd.cmd.Start(); err != nil {
		return err
	}

	if timeout > 0 {
		go func() {
			timer := time.NewTimer(time.Second * time.Duration(timeout))
			select {
			case <-cmd.done:
				timer.Stop()
			case <-timer.C:
				_ = cmd.Kill()
			}
		}()
	}

	go func() {
		var err error
		defer func() {
			close(cmd.done)
			cmd.doneFn(cmd, err)
		}()
		err = cmd.cmd.Wait()
	}()
	return nil
}

/* 以下方法需Run调用成功后才能使用 */

// 判断程序结束
func (cmd *Cmd) Done() bool {
	select {
	case <-cmd.done:
		return true
	default:
		return false
	}
}

func (cmd *Cmd) Kill() error {
	if cmd.cmd.Process == nil {
		return errors.New("command: not started")
	}
	return cmd.cmd.Process.Signal(syscall.SIGKILL)
}

func (cmd *Cmd) Signal(sig syscall.Signal) error {
	if cmd.cmd.Process == nil {
		return errors.New("command: not started")
	}
	return cmd.cmd.Process.Signal(sig)
}

func (cmd *Cmd) Pid() int {
	if cmd.cmd.Process == nil {
		// command: not started
		return -1
	}
	return cmd.cmd.Process.Pid
}

func (cmd *Cmd) ProcessState() *os.ProcessState {
	if !cmd.Done() {
		return nil
	}
	return cmd.cmd.ProcessState
}

func CommandWithCmd(cmd *exec.Cmd) *Cmd {
	return &Cmd{
		cmd:  cmd,
		done: make(chan struct{}),
	}
}

func Command(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	return CommandWithCmd(cmd)
}

const defCnt = 10 // 日志缓存行数

type LineBuffer struct {
	cnt  int
	line [][]byte
	buff *bytes.Buffer
}

func NewLineBuffer(line int) *LineBuffer {
	return &LineBuffer{
		cnt:  line,
		line: make([][]byte, 0, line),
		buff: &bytes.Buffer{},
	}
}

func (this *LineBuffer) Write(p []byte) (n int, err error) {
	n, err = this.buff.Write(p)
	if err != nil {
		return
	}
	for {
		line, err := this.buff.ReadBytes('\n')
		if err != nil {
			break
		}
		if len(this.line) > this.cnt {
			this.line = append(this.line[1:], line)
		} else {
			this.line = append(this.line, line)
		}
	}
	return
}

func (this *LineBuffer) Bytes() []byte {
	return bytes.Join(this.line, nil)
}
