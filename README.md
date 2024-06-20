# ReRun

watch a folder and rerun a command whenever files change

go run main.go <watch directory> <run command> <kill ports>"
go run main.go ../Hello \"go run ../Hello/main.go\" 8080,3000"

rerun <watch directory> <run command> <kill ports>"
rerun ../Hello "go run ../Hello/main.go" 8080,3000"

## Missing

Port number has to be passed to kill the server
cannot get the new process, so cannot kill it directly.
