package logger

import (
	"time"

	"github.com/codeharik/rerun/types"
)

type StdLogSave struct {
	savedOutput map[string][]types.LogEntry
	Group       string
	fn          func(string, func(string)) (n int, err error)
}

func CreateStdOutSave(
	savedOutput map[string][]types.LogEntry,
	fn func(string, func(string)) (n int, err error),
) *StdLogSave {
	return &StdLogSave{
		savedOutput: savedOutput,
		fn:          fn,
	}
}

func (so *StdLogSave) Write(p []byte) (n int, err error) {
	return so.fn(string(p), func(x string) {
		so.savedOutput[so.Group] = append(so.savedOutput[so.Group],
			types.LogEntry{
				Timestamp: time.Now(),
				Log:       x,
			})
	})
}
