package logger

import (
	"log"
	"os"
)

type FileLogger struct {
	log *log.Logger
}

func NewLogger(file *os.File) *FileLogger {
	flogger := &FileLogger{
		log: log.New(file, "", log.Ldate|log.Ltime),
	}
	return flogger
}

func (log *FileLogger) Log(message string) {
	log.log.Println(message)
}
