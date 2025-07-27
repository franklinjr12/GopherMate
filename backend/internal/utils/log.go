package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var logFilePath = "tmp/log.txt"
var logMu sync.Mutex

// Log writes a log entry with the given level and message to both console and file.
// Level should be one of: "DEBUG", "INFO", "WARNING", "ERROR" (case-insensitive).
func Log(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Output to console
	fmt.Print(entry)

	// Output to file (append mode, thread-safe)
	logMu.Lock()
	defer logMu.Unlock()
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(entry)
	}
}

func LogDebug(message string) {
	Log("DEBUG", message)
}

func LogInfo(message string) {
	Log("INFO", message)
}

func LogWarning(message string) {
	Log("WARNING", message)
}

func LogError(message string) {
	Log("ERROR", message)
}
