package helper

import (
	"fmt"
	"os"
	"time"
)

func TickerFunction(done chan os.Signal, t time.Duration, fn func()) {
	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			case s := <-done:
				fmt.Println("Stop Timer")
				done <- s
				return
			}
		}
	}()
}
