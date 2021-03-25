package operators

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"os/exec"
)

func DefaultExec(commands string) (res []byte, err error) {
	return exec.CommandContext(context.Background(), "sh", "-c", commands).Output()
}

func ExecWithStreamOutput(commands string, output chan<- string) (res []byte, err error) {
	//return exec.CommandContext(context.Background(), "sh", "-c", commands).Output()
	cmd := exec.CommandContext(context.Background(), "sh", "-c", commands)
	var stdout, stderr io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	if err = cmd.Start(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			output <- m
			fmt.Println(m)
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			output <- m
			fmt.Println(m)
		}
	}()
	if err = cmd.Wait(); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	return res, nil
}
