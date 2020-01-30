package Persisters

import (
	"fmt"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type consoleLogger struct {
	logUntil loggingC.LogType
	colors   map[loggingC.LogType]ConsoleLoggerColorDelegate
}

// NewConsoleLogger -
func NewConsoleLogger(logUntil loggingC.LogType, colors map[loggingC.LogType]ConsoleLoggerColorDelegate) loggingC.Logger {
	return &consoleLogger{
		logUntil: logUntil,
		colors:   colors,
	}
}

// NewConsoleLoggerDefaultColors -
func NewConsoleLoggerDefaultColors() map[loggingC.LogType]ConsoleLoggerColorDelegate {
	return map[loggingC.LogType]ConsoleLoggerColorDelegate{
		loggingC.TypeDisable: WhiteString,
		loggingC.TypeTrace:   BlackStringYellowBG,
		loggingC.TypePanic:   RedString,
		loggingC.TypeFatal:   RedString,
		loggingC.TypeError:   RedString,
		loggingC.TypeWarning: YellowString,
		loggingC.TypeInfo:    WhiteString,
		loggingC.TypeDebug:   CyanString,
	}
}

// ConsoleLoggerColorDelegate -
type ConsoleLoggerColorDelegate func(string, ...interface{}) string

// BlackStringYellowBG -
func BlackStringYellowBG(format string, a ...interface{}) string {
	c := New(FgBlack, BgYellow)
	return c.Sprintf(format, a...)
}

// BlackStringWhiteBG -
func BlackStringWhiteBG(format string, a ...interface{}) string {
	c := New(FgBlack, BgWhite)
	return c.Sprintf(format, a...)
}

func (thisRef consoleLogger) Log(logEntry loggingC.LogEntry) {
	if logEntry.Type == loggingC.TypeDisable {
		return
	}

	if logEntry.Type > thisRef.logUntil &&
		logEntry.Type != loggingC.TypeTrace {
		return
	}

	if logEntry.Type == loggingC.TypeTrace {
		fmt.Println(thisRef.colors[loggingC.TypeTrace](logEntry.Message))

	} else if logEntry.Type < loggingC.TypeWarning {
		fmt.Println(thisRef.colors[loggingC.TypeError](logEntry.Message))

	} else if logEntry.Type == loggingC.TypeWarning {
		fmt.Println(thisRef.colors[loggingC.TypeWarning](logEntry.Message))

	} else if logEntry.Type == loggingC.TypeInfo {
		fmt.Println(thisRef.colors[loggingC.TypeInfo](logEntry.Message))

	} else if logEntry.Type == loggingC.TypeDebug {
		fmt.Println(thisRef.colors[loggingC.TypeDebug](logEntry.Message))
	}
}
