# ReRun

* `Monitor a Directory`: Watch a specified directory for any file changes.
* `Automatic Command Execution`: Run a command automatically when files in the directory change.
* `Periodic Rerun`: Option to rerun the command at regular intervals.
* `Customizable`: Set specific ports to kill if needed.

```go
ReRun: Monitor a directory and automatically execute a command when files change, or rerun the command on a set interval.
  -k string
        Optional Kill Ports
  -t int
        Optional Rerun Delay Time in Milliseconds[Min 100] (default -1)

Usage: go run main.go [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>
Usage: go run main.go ../Hello "go run ../Hello/main.go"
Usage: go run main.go -k=8080,3000 -t=4000 ../Hello "go run ../Hello/main.go"

Usage: rerun [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>
Usage: rerun -k=8080,3000 -t=4000 ../Hello "go run ../Hello/main.go"
Usage: rerun ../Hello "go run ../Hello/main.go"
```
## Installation

```bash
go get github.com/codeharik/rerun
go install github.com/codeharik/rerun
```
