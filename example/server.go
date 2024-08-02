package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)

		fmt.Println("GET /")
		fmt.Fprintln(w, `
			<body style="background:#e9e0ff;color:purple;text-align: center;align-content: center;font: 30px monospace;">
				<a href="/docs">Go to Docs</a>
			</body>`)
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)

		fmt.Println("GET /docs")
		fmt.Fprintln(w, `
			<body style="background:white;color:purple;text-align: center;align-content: center;font: 30px monospace;">
				<a href="/">Go to Home</a>
			</body>`)
	})

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
