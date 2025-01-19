// logger/logger.go

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
	//err     error
)

func Init() {
	// Get the absolute path to the project root
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get project root:", err)
		os.Exit(1)
	}

	// Construct the log file path within the project root
	logFilePath := filepath.Join(projectRoot, "application.log")

	// Create or open the log file
	logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
}

func Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	log.Fatalf("[FATAL] "+msg, args...)
}

func LogWithTimestamp(level string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [%s] %s", timestamp, level, fmt.Sprint(args...))
}
