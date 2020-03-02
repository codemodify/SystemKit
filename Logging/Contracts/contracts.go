package Contracts

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogType - Log type
type LogType int

// LogType -
const (
	TypeDisable LogType = iota // 0
	TypeTrace                  // 1 - log this no matter what
	TypePanic                  // 2
	TypeFatal                  // 3
	TypeError                  // 4
	TypeWarning                // 5
	TypeInfo                   // 6
	TypeSuccess                // 7
	TypeDebug                  // 8
)

func (logType LogType) String() string {
	switch logType {
	case TypeDisable:
		return "Disable"
	case TypeTrace:
		return "TRACE"
	case TypePanic:
		return "PANIC"
	case TypeFatal:
		return "FATAL"
	case TypeError:
		return "ERROR"
	case TypeWarning:
		return "WAARN"
	case TypeInfo:
		return "INFOO"
	case TypeSuccess:
		return "SUCES"
	case TypeDebug:
		return "DEBUG"

	default:
		return fmt.Sprintf("%d", int(logType))
	}
}

// Fields -
type Fields map[string]interface{}

func (thisRef Fields) String() string {
	bytes, err := json.Marshal(thisRef)
	if err != nil {
		return ""
	}

	return string(bytes)
}

// LogEntry -
type LogEntry struct {
	Time    time.Time // time.Now()
	Type    LogType   // TypeDisable .. -> .. TypeDebug
	Tag     string    // "This-Is-A-Tag"
	Level   int       // Ex: parentMethod - level 0, childMethod() - level 1, useful for concurrent sorted logging with call-stack alike
	Message string    //
}

// Logger -
type Logger interface {
	Log(logEntry LogEntry)
}

// EasyLogger -
type EasyLogger interface {
	Logger

	KeepOnlyLogs(logTypes []LogType)

	LogPanicWithTagAndLevel(tag string, level int, message string)
	LogFatalWithTagAndLevel(tag string, level int, message string)
	LogErrorWithTagAndLevel(tag string, level int, message string)
	LogWarningWithTagAndLevel(tag string, level int, message string)
	LogInfoWithTagAndLevel(tag string, level int, message string)
	LogSuccessWithTagAndLevel(tag string, level int, message string)
	LogDebugWithTagAndLevel(tag string, level int, message string)
	LogTraceWithTagAndLevel(tag string, level int, message string)

	LogPanic(message string)
	LogFatal(message string)
	LogError(message string)
	LogWarning(message string)
	LogInfo(message string)
	LogSuccess(message string)
	LogDebug(message string)
	LogTrace(message string)

	LogPanicWithFields(fields Fields)
	LogFatalWithFields(fields Fields)
	LogErrorWithFields(fields Fields)
	LogWarningWithFields(fields Fields)
	LogInfoWithFields(fields Fields)
	LogSuccessWithFields(fields Fields)
	LogDebugWithFields(fields Fields)
	LogTraceWithFields(fields Fields)
}
