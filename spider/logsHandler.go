package spider

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/codeharik/rerun/types"
)

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
