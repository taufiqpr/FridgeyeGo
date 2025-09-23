package helper

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger

func init() {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0755)
	}
	logFile := filepath.Join("log", time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
}

func Info(msg string)  { logger.Println("[INFO] " + msg) }
func Error(msg string) { logger.Println("[ERROR] " + msg) }
