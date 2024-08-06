package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	version := flag.String("v", "default_version", "Version number")
	flag.Parse()

	fmt.Println(*version)

	saveToFile(
		"info.go",
		strings.ReplaceAll(info, "${VERSION}", *version),
	)

	readmeTmpl, err := os.ReadFile("tools/tmpl.README")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	out, err := exec.Command("sh", "-c", "go run .").CombinedOutput()
	if err != nil {
		log.Fatalf(" %v", err)
	} else {
		saveToFile("README.md", strings.ReplaceAll(string(readmeTmpl), "${MESSAGE}", string(out)))
	}

	fmt.Println("Code generation complete.")
}

func saveToFile(fileName string, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file: %v", err)
	}
}

const info = `
package main

const Version = "${VERSION}"
`
