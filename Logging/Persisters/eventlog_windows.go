// +build windows

package Persisters

import (
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type windowsEventlogLogger struct {
	logUntil loggingC.LogType
}

// NewWindowsEventLogLogger -
func NewWindowsEventlogLogger(logUntil loggingC.LogType) loggingC.Logger {
	return &windowsEventlogLogger{
		logUntil: logUntil,
	}
} 

func (thisRef windowsEventlogLogger) Log(logEntry loggingC.LogEntry) {
	wel, err := eventlog.Open(logEntry.Tag)
	if err != nil {
		return
	}
	defer elog.Close()

	if logEntry.Type == loggingC.TypeDisable {
		return
	}

	if logEntry.Type > thisRef.logUntil &&
		logEntry.Type != loggingC.TypeTrace {
		return
	}

	if logEntry.Type < loggingC.TypeWarning {
		fmt.Println(RedString(logEntry.Message))
		wel.Err(1, logEntry.Message)
	} else if logEntry.Type == loggingC.TypeWarning {
		wel.Warn(1, logEntry.Message)
	} else if logEntry.Type == loggingC.TypeInfo {
		wel.Info(1, logEntry.Message)
	} else if logEntry.Type == loggingC.TypeDebug {
		wel.Debug(1, logEntry.Message)
	}
}
