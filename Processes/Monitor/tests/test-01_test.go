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

func Test_01(t *testing.T) {
	logging.Init(logging.NewEasyLoggerForLogger(loggingP.NewFileLogger(loggingC.TypeDebug, "log1.log")))

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
		// Executable: "htop",
		Executable: "sh",
		Args:       []string{"-c", "while :; do echo 'Hit CTRL+C'; echo aaaaaaa 1>&2; sleep 1; done"},
		OnStdOut: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("OnStdOut: %v", string(data)),
			})
		},
		OnStdErr: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("OnStdErr: %v", string(data)),
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

	time.Sleep(5 * time.Second)

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
