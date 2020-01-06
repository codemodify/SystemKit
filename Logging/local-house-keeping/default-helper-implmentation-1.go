package housekeeping

import (
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

type defaultHelperImplmentation struct {
	loggerToSendTo loggingC.Logger
	logTypes       []loggingC.LogType
}

// NewDefaultHelperImplmentation -
func NewDefaultHelperImplmentation(logger loggingC.Logger) loggingC.EasyLogger {
	return &defaultHelperImplmentation{
		loggerToSendTo: logger,
		logTypes: []loggingC.LogType{
			loggingC.TypePanic,
			loggingC.TypeFatal,
			loggingC.TypeError,
			loggingC.TypeWarning,
			loggingC.TypeInfo,
			loggingC.TypeDebug,
		},
	}
}

func (thisRef *defaultHelperImplmentation) KeepOnlyLogs(logTypes []loggingC.LogType) {
	thisRef.logTypes = logTypes
}

func (thisRef defaultHelperImplmentation) Log(logEntry loggingC.LogEntry) {
	thisRef.loggerToSendTo.Log(logEntry)
}
