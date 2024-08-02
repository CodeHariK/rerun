package spider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
}
