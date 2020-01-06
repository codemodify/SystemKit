package tests

import (
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingF "github.com/codemodify/SystemKit/Logging/Formatters"
	loggingM "github.com/codemodify/SystemKit/Logging/Mixers"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"
)

func Test_03(t *testing.T) {
	logging.Init(
		loggingF.NewSimpleFormatterLogger(
			loggingM.NewAsyncLogger(
				loggingM.NewBufferedLogger(
					loggingM.NewMultiLogger(
						[]loggingC.Logger{
							loggingP.NewConsoleLogger(loggingC.TypeDebug),
							loggingP.NewFileLogger(loggingC.TypeDebug, "log.log"),
						},
					),
					loggingM.BufferedLoggerConfig{
						MaxLogEntries: 100,
					},
				),
			),
		),
	)

	logging.Instance().LogInfo("Info line")
	logging.Instance().LogWarning("Warning line")
	logging.Instance().LogError("Error line")
}
