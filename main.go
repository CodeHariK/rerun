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

const version = "v0.1.3"

//go:embed spider.html
var spiderhtml embed.FS

func main() {
	flagKillPorts := flag.String("k", "", "Optional Kill Ports")
	flagReRunDelay := flag.Int("t", -1, "Optional Rerun Delay Time in Milliseconds [Min 100]")

	flag.Parse()

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) < 2 {
		fmt.Printf("ReRun %s : Monitor a directory and automatically execute a command when directory change, or rerun the command on a set interval.\n", version)
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("SPIDER : http://localhost:9753")
		fmt.Println()
		fmt.Println("Usage: go run main.go [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: go run main.go example \"go run example/server.go\"")
		fmt.Println("Usage: go run main.go -k=8080,3000 -t=4000 example \"go run example/server.go\"")
		fmt.Println()
		fmt.Println("Usage: rerun [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: rerun example \"go run example/server.go\"")
		fmt.Println("Usage: rerun -k=8080,3000 -t=4000 example \"go run example/server.go\"")
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

	rerun := types.ReRun{}

	var wg sync.WaitGroup
	defer wg.Wait()

	stdOutLogs := make(map[string][]types.LogEntry)
	stdErrLogs := make(map[string][]types.LogEntry)

	spider := spider.NewSpider(directory, spiderhtml, stdOutLogs, stdErrLogs)
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
