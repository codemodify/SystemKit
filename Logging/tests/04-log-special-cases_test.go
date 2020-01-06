package tests

import (
	"testing"
	"time"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingF "github.com/codemodify/SystemKit/Logging/Formatters"
	loggingM "github.com/codemodify/SystemKit/Logging/Mixers"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"
)

// Special Cases Logger - `NewInMemoryGroupedAndSortedLogger()` is super useful in
// concurrent apps for debugging exact sequences of logs occured with-in the tree
// of parent-child threads. The `logLevel` will get increased for the child and also affects
// the tabulation when reading the logs. Looks like a "call-stack" type of thing clearly showing what
// what log line was called from what parent.
//
//		This logger stores the log entries in-memory.
//		Groups the log entries by the LOG-TAG
//		Sorts the log entrie in each group by the time
//		When you call `Flush()` it will save to the persisters
//		a nice-call-stack-alike log lines

func Test_04(t *testing.T) {
	var inMemoryGroupedAndSortedLogger = loggingM.NewInMemoryGroupedAndSortedLogger(
		loggingF.NewSimpleFormatterLogger(
			loggingP.NewFileLogger(loggingC.TypeDebug, "log.log"),
		),
	)

	logging.Init(
		logging.NewEasyLoggerForLogger(inMemoryGroupedAndSortedLogger),
	)

	//
	// ... smart concurrent stuff that gets logged, grouped and sorted by time and call-stack
	//
	go func(logTag string, logLevel int) {
		logging.Instance().LogInfoWithTagAndLevel(logTag, logLevel, "Info line")

		go func(logTag string, logLevel int) {
			logging.Instance().LogInfoWithTagAndLevel(logTag, logLevel, "Info line")

			go func(logTag string, logLevel int) {
				logging.Instance().LogInfoWithTagAndLevel(logTag, logLevel, "Info line")
			}(logTag, logLevel+1)
		}(logTag, logLevel+1)

		logging.Instance().LogInfoWithTagAndLevel(logTag, logLevel, "Info line")
	}("LOG-TAG-1", 0)

	go func(logTag string, logLevel int) {
		logging.Instance().LogWarningWithTagAndLevel(logTag, logLevel, "Warning line")

		go func(logTag string, logLevel int) {
			logging.Instance().LogWarningWithTagAndLevel(logTag, logLevel, "Warning line")

			go func(logTag string, logLevel int) {
				logging.Instance().LogWarningWithTagAndLevel(logTag, logLevel, "Warning line")
			}(logTag, logLevel+1)
		}(logTag, logLevel+1)

		logging.Instance().LogWarningWithTagAndLevel(logTag, logLevel, "Warning line")
	}("LOG-TAG-2", 0)

	go func(logTag string, logLevel int) {
		logging.Instance().LogErrorWithTagAndLevel(logTag, logLevel, "Error line")

		go func(logTag string, logLevel int) {
			logging.Instance().LogErrorWithTagAndLevel(logTag, logLevel, "Error line")

			go func(logTag string, logLevel int) {
				logging.Instance().LogErrorWithTagAndLevel(logTag, logLevel, "Error line")
			}(logTag, logLevel+1)
		}(logTag, logLevel+1)

		logging.Instance().LogErrorWithTagAndLevel(logTag, logLevel, "Error line")
	}("LOG-TAG-3", 0)

	//
	// ... more smart concurrent stuff
	//

	// Dump log entries from memory to the configured pipe-line
	time.Sleep(5 * time.Second)
	inMemoryGroupedAndSortedLogger.Flush()
}
