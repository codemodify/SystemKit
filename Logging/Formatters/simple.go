package Formatters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	housekeeping "github.com/codemodify/SystemKit/Logging/local-house-keeping"
)

type simpleFormatterLogger struct {
	loggerToSendTo loggingC.Logger
}

// NewSimpleFormatterLogger -
func NewSimpleFormatterLogger(logger loggingC.Logger) loggingC.EasyLogger {
	var formatterLogger = &simpleFormatterLogger{
		loggerToSendTo: logger,
	}

	return housekeeping.NewDefaultHelperImplmentation(
		formatterLogger,
	)
}

func (thisRef simpleFormatterLogger) Log(logEntry loggingC.LogEntry) {
	var formattedTime = logEntry.Time.UTC().Format(time.RFC3339Nano)
	var formattedTimeLen = len(formattedTime)
	if formattedTimeLen < 30 {
		var spacesCount = 30 - formattedTimeLen

		var newV = fmt.Sprintf("%"+strconv.Itoa(spacesCount+1)+"v", "Z")
		newV = strings.Replace(newV, " ", "0", spacesCount)

		formattedTime = strings.Replace(
			formattedTime,
			"Z",
			newV,
			1,
		)
	}

	var formatting = "%s | %s"
	if len(strings.TrimSpace(logEntry.Tag)) > 0 {
		formatting = formatting + " | %s"
	} else {
		formatting = formatting + " |"
	}

	if logEntry.Level > 0 {
		formatting = formatting + fmt.Sprintf(" %"+strconv.Itoa(logEntry.Level*4)+"v", "")
		formatting += " ->"
	}
	formatting = formatting + " %s"

	if len(strings.TrimSpace(logEntry.Tag)) > 0 {
		logEntry.Message = fmt.Sprintf(
			formatting,
			formattedTime,
			logEntry.Type,
			logEntry.Tag,
			logEntry.Message,
		)
	} else {
		logEntry.Message = fmt.Sprintf(
			formatting,
			formattedTime,
			logEntry.Type,
			logEntry.Message,
		)
	}

	thisRef.loggerToSendTo.Log(logEntry)
}
