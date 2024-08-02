package spider

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *Spider) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the "/redirect" prefix from the request URI
	targetURI := strings.TrimPrefix(r.RequestURI, "/redirect")

	// Create the URL to forward the request to
	targetURL := fmt.Sprintf("%s%d%s", "http://localhost:", s.watchPort, targetURI)

	// Create a new request to forward
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code and copy the response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
