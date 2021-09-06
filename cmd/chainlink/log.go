package main

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	color string
}

// NewLogger will create a new *log.Logger from the stdlib log package but with custom configurations.
func NewLogger(prefix, color string) *Logger {
	logger := &Logger{
		Logger: log.New(os.Stdout, color+prefix+" | ", log.Ldate|log.Lmicroseconds|log.Lshortfile),
		color:  color,
	}
	return logger
}
