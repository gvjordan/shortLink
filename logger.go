package main

import (
	"log"
	"os"
)

type logger struct {
	PrintWarning *log.Logger
	PrintError   *log.Logger
	PrintInfo    *log.Logger
	PrintDebug   *log.Logger
}

func newLogger(path string) *logger {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	return &logger{
		PrintWarning: log.New(f, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		PrintError:   log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		PrintInfo:    log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		PrintDebug:   log.New(f, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// logging levels:
// 		 0: error
// 		 1: warning
// 		 2: info
// 		 3: debug

func (l *logger) Error(msg string) {
	if c.Logs && c.LogsLevel >= 0 {
		l.PrintError.Println(msg)
	}
}

func (l *logger) Warning(msg string) {
	if c.Logs && c.LogsLevel >= 1 {
		l.PrintWarning.Println(msg)
	}
}

func (l *logger) Info(msg string) {
	if c.Logs && c.LogsLevel >= 2 {
		l.PrintInfo.Println(msg)
	}
}

func (l *logger) Debug(msg string) {
	if c.Logs && c.LogsLevel >= 3 {
		l.PrintInfo.Println(msg)
	}
}
