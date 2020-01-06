package Logging

import (
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"
	housekeeping "github.com/codemodify/SystemKit/Logging/local-house-keeping"
)

var instance loggingC.EasyLogger

// Instance -
func Instance() loggingC.EasyLogger {
	return instance
}

// Init -
func Init(logger loggingC.EasyLogger) {
	instance = logger
}

// NewConsoleLogger -
func NewConsoleLogger() loggingC.EasyLogger {
	return housekeeping.NewDefaultHelperImplmentation(
		loggingP.NewConsoleLogger(loggingC.TypeDebug),
	)
}

// NewFileLogger -
func NewFileLogger() loggingC.EasyLogger {
	return housekeeping.NewDefaultHelperImplmentation(
		loggingP.NewFileLoggerDefaultName(loggingC.TypeDebug),
	)
}

func NewEasyLoggerForLogger(logger loggingC.Logger) loggingC.EasyLogger {
	return housekeeping.NewDefaultHelperImplmentation(logger)
}
