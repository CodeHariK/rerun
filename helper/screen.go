package helper

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func ClearScreen() {
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd := exec.Command("clear") // Use "clear" command for macOS and Linux
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") // Use "cls" command for Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported operating system.")
	}
}

func Spinner(t time.Duration) {
	timeUp := time.After(t)
	for {
		select {
		case <-timeUp:
			return
		default:
			for _, c := range `....` {
				fmt.Printf("\r%c", c)
				time.Sleep(time.Millisecond * 30)
			}
		}
	}
}
