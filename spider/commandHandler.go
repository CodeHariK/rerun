package spider

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/codeharik/rerun/logger"
	"github.com/codeharik/rerun/types"
)

func (s *Spider) executeCommandWithContext(ctx context.Context, command *exec.Cmd) {
	defer s.wg.Done()

	fmt.Println("Executing command:", command.Args)

	if err := command.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	done := make(chan error)
	go func() {
		fmt.Println("---NOT-Done----")
		done <- command.Wait()
		fmt.Println("---Done----")
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Command cancelled")
		if err := command.Process.Kill(); err != nil {
			fmt.Printf("Error killing command: %v\n", err)
		}
	case err := <-done:
		if err != nil {
			fmt.Printf("Command execution failed: %v\n", err)
		} else {
			fmt.Println("Command executed successfully")
		}
	}

	fmt.Println("~~~~~~~~~~")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.runningCommand = nil
	s.cancelFunc = nil
}

func (s *Spider) handleExecute(command ...string) {
	stdOutLogs := logger.CreateStdOutSave(
		make(map[string][]types.LogEntry),
		func(p string, append func(string)) (n int, err error) {
			s.BroadcastMessage(fmt.Sprintf("Console:Output:%s", string(p)), Connection{ID: "SPIDER"})
			// return os.Stdout.Write(p)
			return len(p), nil
		},
	)

	stdErrLogs := logger.CreateStdOutSave(
		make(map[string][]types.LogEntry),
		func(p string, append func(string)) (n int, err error) {
			s.BroadcastMessage(fmt.Sprintf("Console:Error:%s", string(p)), Connection{ID: "SPIDER"})
			// return os.Stderr.Write(p)
			return len(p), nil
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", append([]string{"-c"}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdOutLogs
	cmd.Stderr = stdErrLogs

	// Cancel any existing command before starting a new one
	if s.cancelFunc != nil {
		fmt.Println("Cancelling previous running command")
		s.cancelFunc()
		s.wg.Wait()
		fmt.Println("Cancelled previous running command")
	}

	s.mu.Lock()
	s.runningCommand = cmd
	s.cancelFunc = cancel
	s.mu.Unlock()

	s.wg.Add(1)
	fmt.Println("Going to start")
	go s.executeCommandWithContext(ctx, cmd)
}

// func (s *Spider) handleExecute(command ...string) {
// 	go func() {
// 		stdOutLogs := logger.CreateStdOutSave(
// 			make(map[string][]string),
// 			func(p []byte) (n int, err error) {
// 				s.BroadcastMessage(fmt.Sprintf("Console Output %s", string(p)), Connection{ID: "SPIDER"})
// 				// return os.Stdout.Write(p)
// 				return len(p), nil
// 			},
// 		)

// 		stdErrLogs := logger.CreateStdOutSave(
// 			make(map[string][]string),
// 			func(p []byte) (n int, err error) {
// 				s.BroadcastMessage(fmt.Sprintf("Console Error %s", string(p)), Connection{ID: "SPIDER"})
// 				// return os.Stderr.Write(p)
// 				return len(p), nil
// 			},
// 		)

// 		// Execute the command
// 		s.runningCommand = exec.Command("sh", append([]string{"-c"}, command...)...)
// 		s.runningCommand.Stdin = os.Stdin
// 		s.runningCommand.Stdout = stdOutLogs
// 		s.runningCommand.Stderr = stdErrLogs

// 		// output, err := cmd.CombinedOutput()
// 		// if err != nil {
// 		// 	fmt.Printf("Exec Error : %v", err)
// 		// }

// 		// s.BroadcastMessage(fmt.Sprintf("Console %s", string(output)), Connection{ID: "SPIDER"})

// 		fmt.Println("-----")
// 		fmt.Println("Start")
// 		fmt.Println("-----")
// 		if err := s.runningCommand.Start(); err != nil {
// 			fmt.Println("---------")
// 			fmt.Printf("Error Exec : %v\n", err)
// 			fmt.Println("---------")
// 		}

// 		// Wait for the command to finish
// 		if err := s.runningCommand.Wait(); err != nil {
// 			fmt.Println("---------")
// 			fmt.Printf("Command execution failed: %v\n", err)
// 			fmt.Println("---------")
// 		}

// 		fmt.Println("---------")
// 		fmt.Println("Command execution completed successfully.")
// 		fmt.Println("---------")

// 		s.runningCommand = nil
// 	}()
// }
