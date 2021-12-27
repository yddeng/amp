package client

import (
	"testing"
	"time"
)

func TestNewBash(t *testing.T) {
	b := NewBash("go run ../test/test.go", time.Second*2)
	go func() {
		if err := b.Start(); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	t.Log(b.Pid())

	t.Log(b.StdOut())
	if b.HasErr() {
		t.Log(b.StdErr())
	}

}
