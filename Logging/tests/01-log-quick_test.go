package tests

import (
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
)

func Test_01(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfo("Info line")
	logging.Instance().LogWarning("Warning line")
	logging.Instance().LogError("Error line")
}
