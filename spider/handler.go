package spider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/codeharik/rerun/helper"
	"github.com/codeharik/rerun/logger"
	"github.com/codeharik/rerun/types"
	"github.com/gorilla/websocket"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handlePage)
	mux.HandleFunc("GET /file", getFileHandler)
	mux.HandleFunc("POST /save", saveFileHandler)

	mux.HandleFunc("GET /logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.handleLog(w, r)
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWebSocket(w, r)
	})
	server := http.Server{
		Addr:    ":9753",
		Handler: mux,
	}
	return &server
}

func CombineAndSortLogs(logs1, logs2 []types.LogEntry) []types.LogEntry {
	combinedLogs := append(logs1, logs2...)

	sort.Slice(combinedLogs, func(i, j int) bool {
		return combinedLogs[i].Timestamp.Before(combinedLogs[j].Timestamp)
	})

	return combinedLogs
}

func (s *Spider) handleLog(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	jsonData, err := json.Marshal(CombineAndSortLogs(s.stdOutLogs[idString], s.stdErrLogs[idString]))
	if err != nil {
		http.Error(w, "Error encoding combined logs to JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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

func (s *Spider) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	connection := Connection{
		conn: conn,
		ID:   r.RemoteAddr, // Use the remote address as a simple identifier
	}

	fmt.Printf("Adding %s\n", connection.ID)
	s.addConn <- connection

	go func() {
		defer func() {
			s.removeConn <- connection
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					fmt.Printf("Connection closed : %s %v\n", connection.ID, err)
				}
				break
			}

			command := strings.Split(string(message), ":")
			fmt.Println("*")
			fmt.Println(command)
			fmt.Println(string(message) == "SPIDER:Console:Cancel")
			fmt.Println(s.runningCommand)
			fmt.Println("*")

			if string(message) == "SPIDER:PWD" {
				helper.Pwd(s.directory)
				s.BroadcastMessage(fmt.Sprintf("PWD:%s", helper.Pwd(s.directory)), Connection{ID: "SPIDER"})

			}

			if string(message) == "SPIDER:Console:Cancel" {
				if s.cancelFunc != nil {
					fmt.Println("Cancelling previous running command")
					s.cancelFunc()
					s.wg.Wait()
					fmt.Println("Cancelled previous running command")
				}

				continue
			}

			// Handle New Command
			if len(command) > 2 && command[1] == "Console" {
				s.handleExecute(command[2:]...)
			}

			// // SPIDER:Console:Cancel
			// if string(message) == "SPIDER:Console:Cancel" {
			// 	if s.runningCommand != nil {
			// 		fmt.Println("++++++++++++++++++++")
			// 		fmt.Println("cancelTerminalChan 1")
			// 		fmt.Println("++++++++++++++++++++")

			// 		if err := s.runningCommand.Process.Kill(); err != nil {
			// 			fmt.Printf("Cmd Process Kill : %v\n", err)
			// 		}
			// 	}
			// 	continue
			// }
			// // SPIDER:Console:ping google.com
			// if len(command) > 2 && command[1] == "Console" {
			// 	if s.runningCommand != nil {
			// 		fmt.Println("====================")
			// 		fmt.Println("cancelTerminalChan 2")
			// 		fmt.Println("====================")

			// 		if err := s.runningCommand.Process.Kill(); err != nil {
			// 			fmt.Println("---------")
			// 			fmt.Printf("Cmd Process Kill : %v\n", err)
			// 			fmt.Println("---------")
			// 		}

			// 		if err := s.runningCommand.Wait(); err != nil {
			// 			fmt.Println("---------")
			// 			fmt.Printf("Cmd execution failed: %v\n", err)
			// 			fmt.Println("---------")
			// 		}

			// 		// fmt.Println("=====")
			// 		// fmt.Println("sleep")
			// 		// time.Sleep(1 * time.Second)
			// 		// fmt.Println("=====")
			// 	}
			// 	s.handleExecute(command[2:]...)
			// }

			s.BroadcastMessage(fmt.Sprintf("Message:%s", string(message)), connection)
		}
	}()
}

func (s *Spider) handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	spiderhtmlBytes, err := s.spiderhtml.ReadFile("spider.html")
	if err != nil {
		http.Error(w, "Failed to read openapi.json", http.StatusInternalServerError)
		return
	}
	w.Write(spiderhtmlBytes)
}

type SaveRequest struct {
	CurrentFile string `json:"currentFile"`
	Content     string `json:"content"`
}

func saveFileHandler(w http.ResponseWriter, r *http.Request) {
	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(r)
	fmt.Println(r.Body)
	fmt.Println(req)

	if req.CurrentFile == "" {
		http.Error(w, "File path cannot be empty", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(req.CurrentFile); os.IsNotExist(err) {
		http.Error(w, "File does not exist", http.StatusBadRequest)
		return
	}

	// Write updated content
	err := os.WriteFile(req.CurrentFile, []byte(req.Content), 0o644)
	if err != nil {
		http.Error(w, "Error updating file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Content updated successfully")
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("filepath")
	if filePath == "" {
		http.Error(w, "Missing filepath parameter", http.StatusBadRequest)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")

	// Copy the file content to the response writer
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, fmt.Sprintf("Unable to copy file content to response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

const htmlContent = `

`
