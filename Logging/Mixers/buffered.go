package Mixers

import (
	"sync"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

// BufferedLoggerConfig -
type BufferedLoggerConfig struct {
	MaxLogEntries int
}

type bufferedLogger struct {
	loggerToSendTo loggingC.Logger
	logEntries     []loggingC.LogEntry
	rwMutex        sync.RWMutex
	config         BufferedLoggerConfig
}

// NewBufferedLogger -
func NewBufferedLogger(logger loggingC.Logger, config BufferedLoggerConfig) loggingC.Logger {
	return &bufferedLogger{
		loggerToSendTo: logger,
		logEntries:     []loggingC.LogEntry{},
		rwMutex:        sync.RWMutex{},
		config:         config,
	}
}

func (thisRef *bufferedLogger) Log(logEntry loggingC.LogEntry) {
	thisRef.rwMutex.Lock()
	defer thisRef.rwMutex.Unlock()

	thisRef.logEntries = append(thisRef.logEntries, logEntry)

	if len(thisRef.logEntries) > thisRef.config.MaxLogEntries {
		for _, logEntry := range thisRef.logEntries {
			thisRef.loggerToSendTo.Log(logEntry)
		}

		thisRef.logEntries = []loggingC.LogEntry{}
	}
}
