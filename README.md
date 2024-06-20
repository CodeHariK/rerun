# ReRun

watch a folder and rerun a command whenever files change

```bash
go run main.go <watch directory> <run command> <kill ports>"
go run main.go ../Hello \"go run ../Hello/main.go\" 8080,3000"

rerun <watch directory> <run command> <kill ports>"
rerun ../Hello "go run ../Hello/main.go" 8080,3000"
```

## Installation

```bash
go get github.com/codeharik/rerun
go install github.com/codeharik/rerun
```

## Missing

Port number has to be passed to kill the server
cannot get the new process, so cannot kill it directly.
