package spider

import (
	"fmt"
	"net/http"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rerun", s.handlePage)
	mux.HandleFunc("GET /file", getFileHandler)
	mux.HandleFunc("POST /save", saveFileHandler)

	if s.watchPort != -1 {
		mux.HandleFunc("/redirect/*", s.proxyHandler)
	} else {
		mux.HandleFunc("/redirect/*", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
			<body style="background:black;color:purple;text-align: center;align-content: center;font: 30px monospace;">
				Watch Port Missing
			</body>`)
		})
	}

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

func (s *Spider) handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	data := map[string]any{
		"WatchPort": s.watchPort,
	}

	// Render the template with data
	err := s.spiderhtml.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
