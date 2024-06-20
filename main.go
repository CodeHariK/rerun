package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/CodeHariK/rerun/helper"

	"github.com/fsnotify/fsnotify"
)

var (
	currentCmd *exec.Cmd
	cancelFunc context.CancelFunc
	counter    int32 = 0
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <watch directory> <run command> <kill ports>")
		fmt.Println("Usage: go run main.go ../Hello \"go run ../Hello/main.go\" 8080,3000")
		fmt.Println()
		fmt.Println("Usage: rerun <watch directory> <run command> <kill ports>")
		fmt.Println("Usage: rerun ../Hello \"go run ../Hello/main.go\" 8080,3000")
		return
	}

	directory := os.Args[1]
	command := os.Args[2]
	killPorts := []int{}
	if len(os.Args) >= 3 {
		k, err := helper.ParseStringInts(os.Args[3])
		if err == nil {
			killPorts = k
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		runCommand(command, killPorts)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("modified file:", event.Name)
					runCommand(command, killPorts)
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						err = addRecursive(watcher, event.Name)
						if err != nil {
							log.Println("error adding directory:", err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = addRecursive(watcher, directory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func addRecursive(watcher *fsnotify.Watcher, directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("Watching:", path)
		}
		return nil
	})
}

func portKiller(killPorts []int) {
	for _, i := range killPorts {
		fmt.Printf("Kill any process using port %d\n", i)
		killCmd := exec.CommandContext(context.Background(), "sh", "-c", fmt.Sprintf("lsof -ti tcp:%d | xargs kill -9", i))
		if err := killCmd.Run(); err != nil {
			log.Printf("failed to kill process on port 8080: %v", err)
		}
	}
}

func runCommand(command string, killPort []int) {
	helper.ClearScreen()

	atomic.AddInt32(&counter, 1)
	fmt.Printf("\n%d %s\n\n", atomic.LoadInt32(&counter), command)

	portKiller(killPort)

	// TODO : Can't get process
	// Cancel the previous command if it's still running
	// killProcess()

	// Create a new context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel

	parts := strings.Fields(command)
	fn := parts[0]
	args := parts[1:]
	cmd := exec.CommandContext(ctx, fn, args...)

	// Create pipes to capture the command's stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("failed to create stderr pipe: %v", err)
	}

	fmt.Println("Start the command")
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start command: %v", err)
	}
	currentCmd = cmd

	// Copy the command's stdout and stderr to the current process's stdout and stderr
	go io.Copy(os.Stdout, stdoutPipe)
	go io.Copy(os.Stderr, stderrPipe)

	// Wait for the command to finish
	// go func() {
	// 	err = cmd.Wait()
	// 	if err != nil {
	// 		log.Println("Error running command:", err)
	// 	}
	// }()
}

// func killProcess() {
// 	if currentCmd != nil && currentCmd.Process != nil {
// 		fmt.Println("----")
// 		fmt.Println(currentCmd.Process.Pid)
// 		fmt.Println(currentCmd.Process)
// 		fmt.Println("----")
// 		currentCmd.Process.Signal(syscall.SIGINT)
// 		currentCmd.Process.Signal(os.Interrupt)
// 		currentCmd.Process.Signal(os.Kill)
// 	}

// 	if currentCmd != nil && currentCmd.Process != nil {
// 		log.Printf("Attempting to kill process with PID: %d", currentCmd.Process.Pid)
// 		if err := currentCmd.Process.Kill(); err != nil {
// 			log.Printf("failed to kill process: %v", err)
// 		} else {
// 			log.Printf("Successfully killed process with PID: %d", currentCmd.Process.Pid)
// 		}
// 		if cancelFunc != nil {
// 			cancelFunc()
// 		}
// 		log.Printf("Successfully cancelled process with PID: %d", currentCmd.Process.Pid)

// 		time.Sleep(time.Second * 1)

// 		currentCmd = nil
// 	}
// }
