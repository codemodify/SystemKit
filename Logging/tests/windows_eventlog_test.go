// +build windows

package tests

import (
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
)

func Test_windows_eventlog(t *testing.T) {
	logging.Init(
		logging.NewWindowsEventLogger(),
	)

	logging.Instance().LogInfo("Info line")
	logging.Instance().LogWarning("Warning line")
	logging.Instance().LogError("Error line")
}
