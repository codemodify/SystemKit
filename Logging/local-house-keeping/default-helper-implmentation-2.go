package housekeeping

import (
	"time"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

func (thisRef defaultHelperImplmentation) LogPanicWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypePanic,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogFatalWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeFatal,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogErrorWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeError,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogWarningWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeWarning,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogInfoWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeInfo,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogDebugWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeDebug,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
func (thisRef defaultHelperImplmentation) LogTraceWithTagAndLevel(tag string, level int, message string) {
	thisRef.Log(loggingC.LogEntry{
		Time:    time.Now(),
		Type:    loggingC.TypeTrace,
		Tag:     tag,
		Level:   level,
		Message: message,
	})
}
