package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var Log *log.Logger

var logFile *os.File

func Init(filePath string) error {
	if filePath == "" {
		filePath = "logs/app.log"
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logFile = f
	Log = log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags|log.Lmicroseconds)
	return nil
}

func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func TCPEvent(event string, args ...interface{}) {
	msg := fmt.Sprintf("[TCP] "+event, args...)
	Log.Println(msg)
}
