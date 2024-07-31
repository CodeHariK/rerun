package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type ReRun struct{}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Log       string    `json:"log"`
}

func EncodeLog(log string) string {
	entry := LogEntry{
		Timestamp: time.Now(),
		Log:       log,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return ""
	}

	return string(jsonData)
}
