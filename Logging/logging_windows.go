// +build windows

package Logging

import (
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"
	housekeeping "github.com/codemodify/SystemKit/Logging/local-house-keeping"
)

// NewWindowsEventLogger -
func NewWindowsEventLogger() loggingC.EasyLogger {
	return housekeeping.NewDefaultHelperImplmentation(
		loggingP.NewWindowsEventlogLogger(loggingC.TypeDebug),
	)
}
