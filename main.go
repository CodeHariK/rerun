package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/codeharik/rerun/helper"
	"github.com/codeharik/rerun/spider"
	"github.com/codeharik/rerun/types"
	"github.com/codeharik/rerun/watcher"
)

const version = "v0.1.5"

//go:embed ui
var spiderhtml embed.FS

func main() {
	flagKillPorts := flag.String("k", "", "Optional Kill Ports")
	flagReRunDelay := flag.Int("t", -1, "Optional Rerun Delay Time in seconds [Min 1s]")
	watchPort := flag.Int("w", -1, "Optional Watch port")

	flag.Parse()

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) < 2 {
		fmt.Printf("ReRun %s : Monitor a directory and automatically execute a command when directory change, or rerun the command on a set interval.\n", version)
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("SPIDER : http://localhost:9753/ui")
		fmt.Println()
		fmt.Println("Usage: go run main.go [-w Watch Ports] [-k Kill Ports] [-t Rerun Delay Time] <watch directory> <run command>")
		fmt.Println("Usage: go run main.go -w=8080 example \"go run example/server.go\"")
		fmt.Println("Usage: go run main.go -w=8080 -k=8080,3000 -t=30 example \"go run example/server.go\"")
		fmt.Println()
		fmt.Println("Usage: rerun [-w Watch Ports] [-k Kill Ports] [-t Rerun Delay Time] <watch directory> <run command>")
		fmt.Println("Usage: rerun -w=8080 example \"go run example/server.go\"")
		fmt.Println("Usage: rerun -w=8080 -k=8080,3000 -t=30 example \"go run example/server.go\"")
		return
	}

	killPortsString := *flagKillPorts
	rerunTimer := time.Duration(*flagReRunDelay) * 1000000000
	if rerunTimer > 0 && rerunTimer < time.Second*1 {
		log.Fatal("Min 1 seconds delay required")
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

	rerun := types.ReRun{}

	var wg sync.WaitGroup
	defer wg.Wait()

	stdOutLogs := make(map[string][]types.LogEntry)
	stdErrLogs := make(map[string][]types.LogEntry)

	spider := spider.NewSpider(directory, spiderhtml, *watchPort, stdOutLogs, stdErrLogs)
	spider.StartSpider(&wg)

	w := watcher.NewWatcher(
		&rerun,
		command, rerunTimer, killPorts, directory,
		spider,
		stdOutLogs,
		stdErrLogs,
	)
	w.StartWatcher()
}
