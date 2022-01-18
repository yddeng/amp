package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	ecmd := exec.Command("./test")
	//ecmd.Dir = "./test"

	f, err := os.OpenFile("test_err.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	//errout := bytes.Buffer{}
	ecmd.Stderr = f
	ecmd.Stdout = f
	defer f.Close()

	err = ecmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ecmd.Process.Pid, ecmd.ProcessState)
	time.AfterFunc(time.Second, func() {
		syscall.Kill(ecmd.Process.Pid, syscall.SIGKILL)
	})
	err = ecmd.Wait()

	data, err := ioutil.ReadFile("test_err.log")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("end", err, "-", string(data), "-")
}
