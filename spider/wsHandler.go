package spider

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codeharik/rerun/helper"
	"github.com/gorilla/websocket"
)

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
