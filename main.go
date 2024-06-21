package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/codeharik/rerun/helper"

	"github.com/fsnotify/fsnotify"
)

var (
	shellProcess *os.Process
	childProcess *os.Process
	counter      int32 = 0
)

func main() {
	flagKillPorts := flag.String("k", "", "Optional Kill Ports")
	flagReRunDelay := flag.Int("t", -1, "Optional Rerun Delay Time in Milliseconds[Min 100]")

	flag.Parse()

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) < 2 {
		fmt.Println("ReRun: Monitor a directory and automatically execute a command when files change, or rerun the command on a set interval.")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Usage: go run main.go [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: go run main.go ../Hello \"go run ../Hello/main.go\"")
		fmt.Println("Usage: go run main.go -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
		fmt.Println()
		fmt.Println("Usage: rerun [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: rerun -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
		fmt.Println("Usage: rerun ../Hello \"go run ../Hello/main.go\"")
		return
	}

	killPortsString := *flagKillPorts
	rerunTimer := time.Duration(*flagReRunDelay) * 1000000
	if rerunTimer > 0 && rerunTimer < time.Millisecond*100 {
		log.Fatal("Min 100 milliseconds delay required")
	}

	directory := flag.Arg(0)
	command := flag.Arg(1)
	killPorts := []int{}
	if killPortsString != "" {
		k, err := helper.ParseStringInts(killPortsString)
		if err == nil {
			killPorts = k
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan os.Signal)
	defer close(done)
	signal.Notify(done, syscall.SIGINT, os.Kill, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	if rerunTimer >= time.Millisecond*100 {
		helper.TickerFunction(
			done,
			rerunTimer,
			func() {
				runCommand(command, killPorts, rerunTimer)
			},
		)
	}

	go func() {
		defer wg.Done()
		runCommand(command, killPorts, rerunTimer)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("modified file:", event.Name)
					runCommand(command, killPorts, rerunTimer)
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						err = helper.AddRecursive(watcher, event.Name)
						if err != nil {
							fmt.Println("error adding directory:", err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			case <-done:
				if err := watcher.Close(); err != nil {
					fmt.Println("Failed to stop watcher")
				}
				return
			}
		}
	}()

	err = helper.AddRecursive(watcher, directory)
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(command string, killPort []int, t time.Duration) {
	helper.ClearScreen()

	atomic.AddInt32(&counter, 1)
	fmt.Printf("\n%d %s [Rerun:%s]\n\n", atomic.LoadInt32(&counter), command, t)

	helper.KillProcess(shellProcess)
	helper.KillProcess(childProcess)
	helper.PortKiller(killPort)

	cmd := helper.ExecCommand(command)

	helper.Spinner(time.Millisecond * 400)

	helper.CopyProcess(cmd, &shellProcess, &childProcess)

	fmt.Printf("...\n\n")
}
