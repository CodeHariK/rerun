package spider

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"text/template"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/", s.handlePage)
	mux.HandleFunc("GET /file", getFileHandler)
	mux.HandleFunc("POST /save", saveFileHandler)
	mux.HandleFunc("GET /npm/", func(w http.ResponseWriter, r *http.Request) {
		newURL := strings.Replace(r.URL.Path, "/npm/", "https://cdn.jsdelivr.net/npm/", 1)
		http.Redirect(w, r, newURL, http.StatusMovedPermanently)
	})

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
	data := map[string]any{
		"WatchPort": s.watchPort,
	}

	url := strings.TrimPrefix(r.URL.String(), "/")

	if url == "ui/" {
		url = "ui/spider.html"
	}

	if strings.HasSuffix(url, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else {
		w.Header().Set("Content-Type", "text/html")
	}

	tmpl, err := GetTemplate(s.spiderhtml, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template with data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type templateCache struct {
	tmpl *template.Template
	err  error
}

var (
	l     sync.RWMutex
	cache map[string]templateCache = map[string]templateCache{}
)

func ParseTemplate(templatesFS embed.FS, name string) (*template.Template, error) {
	tmpl, err := template.ParseFS(templatesFS, name)
	l.Lock()
	cache[name] = templateCache{tmpl: tmpl, err: err}
	l.Unlock()
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", name, err)
	}
	return tmpl, err
}

func GetTemplate(templatesFS embed.FS, name string) (*template.Template, error) {
	if cacheEntry, ok := cache[name]; ok {
		if cacheEntry.err != nil {
			return nil, cacheEntry.err
		}
		return cacheEntry.tmpl, nil
	}
	return ParseTemplate(templatesFS, name)
}
