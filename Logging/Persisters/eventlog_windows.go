// +build windows

package Persisters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type windowsEventlogLogger struct {
	logUntil        loggingC.LogType
	eventlogLogger  *eventlog.Log
	emergencyLogger *debug.ConsoleLog
}

// NewWindowsEventlogLogger -
func NewWindowsEventlogLogger(logUntil loggingC.LogType) loggingC.Logger {
	binaryName := filepath.Base(os.Args[0])
	emergencyLogger := debug.New(binaryName)

	// _ = eventlog.Remove(binaryName) eventlog.Error | eventlog.Warning | eventlog.Info
	err := eventlog.InstallAsEventCreate(binaryName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {

		if strings.Contains(err.Error(), "registry key already exists") {
			// SAFE to ignore
			// emergencyLogger.Error(1, fmt.Sprint("warning creating service logs: ", err))
		} else if err == windows.ERROR_ACCESS_DENIED {
			// SAFE to ignore
			// most probably running as user
		} else {
			emergencyLogger.Error(1, fmt.Sprint("error creating service logs: ", err))

			return &windowsEventlogLogger{
				logUntil:        logUntil,
				eventlogLogger:  nil,
				emergencyLogger: emergencyLogger,
			}
		}
	}

	eventlogLogger, err := eventlog.Open(binaryName)
	if err != nil {
		emergencyLogger.Error(1, err.Error())

		return &windowsEventlogLogger{
			logUntil:        logUntil,
			eventlogLogger:  nil,
			emergencyLogger: emergencyLogger,
		}
	}

	return &windowsEventlogLogger{
		logUntil:        logUntil,
		eventlogLogger:  eventlogLogger,
		emergencyLogger: nil,
	}
}

func (thisRef windowsEventlogLogger) Log(logEntry loggingC.LogEntry) {
	if logEntry.Type == loggingC.TypeDisable {
		return
	}

	if logEntry.Type > thisRef.logUntil &&
		logEntry.Type != loggingC.TypeTrace {
		return
	}

	if logEntry.Type < loggingC.TypeWarning {
		if thisRef.eventlogLogger != nil {
			thisRef.eventlogLogger.Error(1, logEntry.Message)
		} else {
			thisRef.emergencyLogger.Error(1, logEntry.Message)
		}
	} else if logEntry.Type == loggingC.TypeWarning {
		if thisRef.eventlogLogger != nil {
			thisRef.eventlogLogger.Warning(1, logEntry.Message)
		} else {
			thisRef.emergencyLogger.Warning(1, logEntry.Message)
		}
	} else if logEntry.Type == loggingC.TypeInfo {
		if thisRef.eventlogLogger != nil {
			thisRef.eventlogLogger.Info(1, logEntry.Message)
		} else {
			thisRef.emergencyLogger.Info(1, logEntry.Message)
		}
	} else if logEntry.Type == loggingC.TypeSuccess {
		if thisRef.eventlogLogger != nil {
			thisRef.eventlogLogger.Info(1, logEntry.Message)
		} else {
			thisRef.emergencyLogger.Info(1, logEntry.Message)
		}
	} else if logEntry.Type == loggingC.TypeDebug {
		if thisRef.eventlogLogger != nil {
			thisRef.eventlogLogger.Info(1, logEntry.Message)
		} else {
			thisRef.emergencyLogger.Info(1, logEntry.Message)
		}
	}
}
