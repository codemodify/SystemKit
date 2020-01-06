package Persisters

import (
	"fmt"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type consoleLogger struct {
	logUntil loggingC.LogType
}

// NewConsoleLogger -
func NewConsoleLogger(logUntil loggingC.LogType) loggingC.Logger {
	return &consoleLogger{
		logUntil: logUntil,
	}
}

func (thisRef consoleLogger) Log(logEntry loggingC.LogEntry) {
	if logEntry.Type == loggingC.TypeDisable {
		return
	}

	if logEntry.Type > thisRef.logUntil &&
		logEntry.Type != loggingC.TypeTrace {
		return
	}

	if logEntry.Type < loggingC.TypeWarning {
		fmt.Println(RedString(logEntry.Message))
	} else if logEntry.Type == loggingC.TypeWarning {
		fmt.Println(MagentaString(logEntry.Message))
	} else if logEntry.Type == loggingC.TypeInfo {
		fmt.Println(logEntry.Message)
	} else if logEntry.Type == loggingC.TypeDebug {
		fmt.Println(GreenString(logEntry.Message))
	}
}
