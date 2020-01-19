package main

import (
	"fmt"
	"testing"
	"time"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	loggingP "github.com/codemodify/SystemKit/Logging/Persisters"

	procMon "github.com/codemodify/SystemKit/Processes/Monitor"
)

func Test_02(t *testing.T) {
	logging.Init(logging.NewEasyLoggerForLogger(loggingP.NewFileLogger(loggingC.TypeDebug, "log2.log")))

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": "Test_01()",
	})

	processID := "test-id"

	monitor := procMon.New()

	// starting
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": "START",
	})

	monitor.Spawn(processID, procMon.Process{
		Executable: "htop",
		OnStdOut: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("OnStdOut: %v", data),
			})
		},
		OnStdErr: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("OnStdErr: %v", data),
			})
		},
	})

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			monitor.GetProcessInfo(processID).IsRunning(),
			monitor.GetProcessInfo(processID).ExitCode(),
			monitor.GetProcessInfo(processID).StartedAt(),
			monitor.GetProcessInfo(processID).StoppedAt(),
		),
	})

	time.Sleep(30 * time.Second)

	// stop
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": "STOP",
	})

	monitor.Stop(processID)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			monitor.GetProcessInfo(processID).IsRunning(),
			monitor.GetProcessInfo(processID).ExitCode(),
			monitor.GetProcessInfo(processID).StartedAt(),
			monitor.GetProcessInfo(processID).StoppedAt(),
		),
	})
}
