package main

import (
	"fmt"
	"testing"
	"time"

	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"
	loggingP "github.com/codemodify/systemkit-logging/persisters"

	procMon "github.com/codemodify/SystemKit/Processes/Monitor"
)

const logID = "TEST"

func Test_02(t *testing.T) {
	logging.Init(logging.NewEasyLoggerForLogger(loggingP.NewFileLogger(loggingC.TypeDebug, "log2.log")))

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": fmt.Sprintf("%s: START", logID),
	})

	processID := "test-id"
	monitor := procMon.New()

	err := monitor.Spawn(processID, procMon.Process{
		Executable: "/usr/local/bin/connectd",
		Args:       []string{"-s", "-p", "bmljb2xhZUByZW1vdGUuaXQ=", "6005CA956F875100B2FCE226D44ABF02B6FA1419", "80:00:00:05:3A:00:36:DC", "T45296", "2", "localhost", "0.0.0.0", "12", "0", "0"},
		OnStdOut: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("%s: OnStdOut: %v", logID, string(data)),
			})
		},
		OnStdErr: func(data []byte) {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("%s: OnStdErr: %v", logID, string(data)),
			})
		},
	})
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"message": fmt.Sprintf("%s: ERROR: %s", logID, err.Error()),
		})

		t.FailNow()
	}

	pi := monitor.GetProcessInfo(processID)
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			logID,
			pi.IsRunning(),
			pi.ExitCode(),
			pi.StartedAt(),
			pi.StoppedAt(),
		),
	})

	timeOut := 30 * time.Second
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		i := 1
		for t := range ticker.C {
			logging.Instance().LogDebugWithFields(loggingC.Fields{
				"message": fmt.Sprintf("%s: WAITING... %d - %v", logID, i, t),
			})
			i++
		}
	}()

	time.Sleep(timeOut)
	ticker.Stop()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"message": fmt.Sprintf("%s: STOP", logID),
	})
	err = monitor.Stop(processID)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"message": fmt.Sprintf("%s: ERROR: %s", logID, err.Error()),
		})

		t.FailNow()
	}

	pi = monitor.GetProcessInfo(processID)
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"message": fmt.Sprintf(
			"%s: IsRunning: %v, ExitCode: %v, StartedAt: %v, StoppedAt: %v",
			logID,
			pi.IsRunning(),
			pi.ExitCode(),
			pi.StartedAt(),
			pi.StoppedAt(),
		),
	})
}
