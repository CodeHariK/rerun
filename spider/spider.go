package spider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/codeharik/rerun/types"
	"github.com/gorilla/websocket"
)

// Connection represents a WebSocket connection.
type Connection struct {
	conn *websocket.Conn
	ID   string
}

type Spider struct {
	directory string

	mu         sync.Mutex
	conns      map[string]Connection
	addConn    chan Connection
	removeConn chan Connection

	runningCommand *exec.Cmd
	cancelFunc     context.CancelFunc
	wg             sync.WaitGroup

	stdOutLogs map[string][]types.LogEntry
	stdErrLogs map[string][]types.LogEntry
}

func NewSpider(
	directory string,
	stdOutLogs map[string][]types.LogEntry,
	stdErrLogs map[string][]types.LogEntry,
) *Spider {
	return &Spider{
		directory: directory,

		conns:      make(map[string]Connection),
		addConn:    make(chan Connection),
		removeConn: make(chan Connection),

		stdOutLogs: stdOutLogs,
		stdErrLogs: stdErrLogs,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Check the origin of the request and decide whether to accept or reject it
		return true
	},
}

func (s *Spider) StartSpider(wg *sync.WaitGroup) {
	server := createServer(s)

	go s.manageConnections()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Spider started on :9753")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
		fmt.Println("Spider stopped")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done

		// Create a context with a timeout to allow the server to shut down gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Spider forced to shutdown: %v\n", err)
		}
	}()
}

func (s *Spider) manageConnections() {
	for {
		select {
		case newConn := <-s.addConn:
			s.mu.Lock()
			s.conns[newConn.ID] = newConn // Add to map
			s.mu.Unlock()
			fmt.Printf("Add new connection %s\n", newConn.ID)
		case removeConn := <-s.removeConn:
			s.mu.Lock()
			if conn, ok := s.conns[removeConn.ID]; ok {
				conn.conn.Close()              // Close WebSocket connection
				delete(s.conns, removeConn.ID) // Remove from map
			}
			s.mu.Unlock()
			fmt.Printf("Remove connection %s\n", removeConn.ID)
		}
	}
}

// BroadcastMessage sends a message to all active connections.
func (s *Spider) BroadcastMessage(message string, sender Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, conn := range s.conns {
		if conn == sender {
			continue
		}
		err := conn.conn.WriteMessage(1, []byte(fmt.Sprintf("%s:%s", sender.ID, message)))
		// fmt.Printf("-> %s:%s\n", sender.ID, string(message))
		if err != nil {
			fmt.Printf("Error broadcasting message to %s: %v\n", conn.ID, err)
			s.removeConn <- conn
		}
	}
}
