package tests

import (
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingF "github.com/codemodify/SystemKit/Logging/Formatters"
	loggingM "github.com/codemodify/SystemKit/Logging/Mixers"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"
)

func Test_02(t *testing.T) {
	logging.Init(
		loggingF.NewSimpleFormatterLogger(
			loggingM.NewMultiLogger(
				[]loggingC.Logger{
					loggingP.NewConsoleLogger(loggingC.TypeDebug, loggingP.NewConsoleLoggerDefaultColors()),
					loggingP.NewFileLogger(loggingC.TypeDebug, "log.log"),
				},
			),
		),
	)

	logging.Instance().LogTrace("Trace line")
	logging.Instance().LogPanic("Panic line")
	logging.Instance().LogFatal("Fatal line")
	logging.Instance().LogError("Error line")
	logging.Instance().LogWarning("Warning line")
	logging.Instance().LogInfo("Info line")
	logging.Instance().LogSuccess("Success line")
	logging.Instance().LogDebug("Debug line")
}
