package watcher

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/codeharik/rerun/helper"
	"github.com/codeharik/rerun/logger"
	"github.com/codeharik/rerun/spider"
	"github.com/codeharik/rerun/types"
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	rerun *types.ReRun

	shellProcess *os.Process
	childProcess *os.Process
	counter      int32

	cmdMutex      sync.Mutex
	command       string
	reRunDuration time.Duration
	killPorts     []int
	directory     string

	spider *spider.Spider

	stdOutLogs logger.StdLogSave
	stdErrLogs logger.StdLogSave
}

func NewWatcher(
	rerun *types.ReRun,

	command string,
	reRunDuration time.Duration,
	killPorts []int,
	directory string,

	spiderServer *spider.Spider,

	stdOutLogs map[string][]types.LogEntry,
	stdErrLogs map[string][]types.LogEntry,
) *watcher {
	return &watcher{
		rerun: rerun,

		stdOutLogs: *logger.CreateStdOutSave(
			stdOutLogs,
			func(s string, append func(string)) (n int, err error) {
				append("Output:" + s)
				spiderServer.BroadcastMessage(fmt.Sprintf("Logs:Output:%s", s), spider.Connection{ID: "SPIDER"})
				return os.Stdout.Write([]byte(s))
			},
		),

		stdErrLogs: *logger.CreateStdOutSave(
			stdErrLogs,
			func(s string, append func(string)) (n int, err error) {
				append("Error:" + s)
				spiderServer.BroadcastMessage(fmt.Sprintf("Logs:Error:%s", s), spider.Connection{ID: "SPIDER"})
				return os.Stderr.Write([]byte(s))
			},
		),

		command:       command,
		reRunDuration: reRunDuration,
		killPorts:     killPorts,
		directory:     directory,

		spider: spiderServer,
	}
}

func (w *watcher) StartWatcher() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, os.Kill, os.Interrupt)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w.runCommand()

	if w.reRunDuration >= time.Millisecond*100 {
		helper.TickerFunction(
			w.reRunDuration,
			func() {
				w.runCommand()
			},
		)
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("Watcher.Events closed:", err)
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("modified file:", event.Name)

					w.runCommand()

					w.spider.BroadcastMessage(fmt.Sprintf("PWD:%s", helper.Pwd(w.directory)), spider.Connection{ID: "SPIDER"})
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
					fmt.Println("Watcher.Errors closed")
					return
				}
				if err != nil {
					fmt.Println("Watcher.Errors : ", err)
				}
			case <-done:
				defer wg.Done()
				w.childProcess.Kill()
				w.shellProcess.Kill()
				time.Sleep(300 * time.Millisecond)
				fmt.Println("Watcher stopped")
				if err := watcher.Close(); err != nil {
					fmt.Println("Failed to stop watcher")
				}
				return
			}
		}
	}()

	err = helper.AddRecursive(watcher, w.directory)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *watcher) runCommand() {
	helper.ClearScreen()

	atomic.AddInt32(&w.counter, 1)

	a := atomic.LoadInt32(&w.counter)

	fmt.Printf("\n%d %s [Rerun:%s] SPIDER: http://localhost:9753/ui\n\n", a, w.command, w.reRunDuration)

	w.spider.BroadcastMessage(fmt.Sprintf("ReRun:%d", a), spider.Connection{ID: "SPIDER"})

	KillProcess(w.shellProcess)
	KillProcess(w.childProcess)
	PortKiller(w.killPorts)

	w.cmdMutex.Lock()
	defer w.cmdMutex.Unlock()

	w.stdOutLogs.Group = strconv.Itoa(int(a))
	w.stdErrLogs.Group = strconv.Itoa(int(a))

	cmd := w.ExecCommand()

	helper.Spinner(time.Millisecond * 400)

	CopyProcess(cmd, &w.shellProcess, &w.childProcess)

	fmt.Printf("...\n\n")
}
