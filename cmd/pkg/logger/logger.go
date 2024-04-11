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

// Log: logs into a file, in this demo i'll use a generic logfile
// in order to show the instructions that are executed by the indexer
func (log *FileLogger) Log(message string) {
	log.log.Println(message)
}
