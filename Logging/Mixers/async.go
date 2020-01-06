package Mixers

import (
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type asyncLogger struct {
	loggerToSendTo loggingC.Logger
}

// NewAsyncLogger -
func NewAsyncLogger(logger loggingC.Logger) loggingC.Logger {
	return &asyncLogger{
		loggerToSendTo: logger,
	}
}

func (thisRef asyncLogger) Log(logEntry loggingC.LogEntry) {
	go func(theLogEntryCopy loggingC.LogEntry) {
		thisRef.loggerToSendTo.Log(theLogEntryCopy)
	}(logEntry)
}
