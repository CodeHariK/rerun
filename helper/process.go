package helper

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func KillProcess(process *os.Process) {
	if process != nil {
		if err := process.Kill(); err != nil {
			fmt.Printf("Failed to kill process: %v", err)
		} else {
			fmt.Printf("Killed process with PID: %d\n", process.Pid)
		}
		process = nil
	}
}

func PortKiller(killPorts []int) {
	for _, i := range killPorts {
		fmt.Printf("Kill any process using port %d\n", i)
		killCmd := exec.CommandContext(context.Background(), "sh", "-c", fmt.Sprintf("lsof -ti tcp:%d | xargs kill -9", i))
		if err := killCmd.Run(); err != nil {
			fmt.Printf("Failed to kill process on port 8080: %v", err)
		}
	}
}

func GetChildPID(parentPID int) (int, error) {
	// Get list of processes
	out, err := exec.Command("ps", "-eo", "pid,ppid").Output()
	if err != nil {
		return 0, err
	}

	// Parse the output to find the child PID
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		ppid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		if ppid == parentPID {
			return pid, nil
		}
	}
	return 0, fmt.Errorf("child process not found ")
}

func CopyProcess(cmd *exec.Cmd, shellProcess **os.Process, childProcess **os.Process) {
	shellPID := cmd.Process.Pid

	serverPID, err := GetChildPID(shellPID)
	if err != nil {
		fmt.Printf("Failed to get child PID: %v\n", err)
	}
	fmt.Printf("ShellPID:%d, ChildPID:%d\n", shellPID, serverPID)
	process, err := os.FindProcess(serverPID)
	if err != nil {
		fmt.Printf("Failed to find process: %v\n", err)
	}
	*shellProcess = cmd.Process
	*childProcess = process
}

// Exec Command
// Copy the command's stdout and stderr to the current process's stdout and stderr
func ExecCommand(command string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", command)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create stdout pipe: %v\n", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to create stderr pipe: %v\n", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v\n", err)
	}

	go io.Copy(os.Stdout, stdoutPipe)
	go io.Copy(os.Stderr, stderrPipe)

	return cmd
}
