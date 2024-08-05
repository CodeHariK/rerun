# ReRun

* `Monitor a Directory`: Watch a specified directory for any file changes.
* `Automatic Command Execution`: Run a command automatically when files in the directory change.
* `Local File Editing`: Provides a built-in web server that serves the files from your monitored directory, allowing direct modifications in the browser.
* `Error and Log Monitoring`: Includes features to monitor errors and logs directly in the browser, helping you track issues and outputs.
* `Periodic Rerun`: Option to rerun the command at regular intervals.
* `Customizable`: Set specific ports to kill if needed.

```go
ReRun v0.1.7 : Monitor a directory and automatically execute a command when directory change, or rerun the command on a set interval.
  -k string
        Optional Kill Ports
  -t int
        Optional Rerun Delay Time in seconds [Min 1s] (default -1)
  -w int
        Optional Watch port (default -1)

SPIDER : http://localhost:9753/ui

Usage: go run main.go [-w Watch Ports] [-k Kill Ports] [-t Rerun Delay Time] <watch directory> <run command>
Usage: go run main.go -w=8080 example "go run example/server.go"
Usage: go run main.go -w=8080 -k=8080,3000 -t=30 example "go run example/server.go"

Usage: rerun [-w Watch Ports] [-k Kill Ports] [-t Rerun Delay Time] <watch directory> <run command>
Usage: rerun -w=8080 example "go run example/server.go"
Usage: rerun -w=8080 -k=8080,3000 -t=30 example "go run example/server.go"
```
## Installation

```bash
go get github.com/codeharik/rerun
go install github.com/codeharik/rerun
```
