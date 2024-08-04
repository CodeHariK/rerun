package helper

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TickerFunction(t time.Duration, fn func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, os.Kill, os.Interrupt)

	// con, cancel := context.WithCancel()

	go func() {
		ticker := time.NewTicker(t)

		for {
			select {
			case <-ticker.C:
				fn()
			// case <-con.Done():
			// 	fmt.Println("Timer stopped")
			// 	ticker.Stop()
			// 	return
			case <-done:
				fmt.Println("Timer stopped")
				ticker.Stop()
				return
			}
		}
	}()
}
