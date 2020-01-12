package Persisters

import (
	"fmt"
	"os"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type fileLogger struct {
	file         *os.File
	errorOccured bool
	logUntil     loggingC.LogType
}

func fileOrFolderExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// NewFileLogger -
func NewFileLogger(logUntil loggingC.LogType, fileName string) loggingC.Logger {
	var f *os.File
	var err error

	if !fileOrFolderExists(fileName) {
		f, err = os.Create(fileName)
	} else {
		f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0660)
	}
	if err != nil {
		// return nil, err
		fmt.Println(err)
	}

	return &fileLogger{
		file:         f,
		errorOccured: (err != nil),
		logUntil:     logUntil,
	}
}

// NewFileLoggerDefaultName -
func NewFileLoggerDefaultName(logUntil loggingC.LogType) loggingC.Logger {
	return NewFileLogger(logUntil, fmt.Sprintf("%s.log", os.Args[0]))
}

func (thisRef fileLogger) Log(logEntry loggingC.LogEntry) {
	if thisRef.errorOccured {
		return
	}

	if logEntry.Type == loggingC.TypeDisable {
		return
	}

	if logEntry.Type > thisRef.logUntil &&
		logEntry.Type != loggingC.TypeTrace {
		return
	}

	thisRef.file.WriteString(logEntry.Message + "\n")
}
