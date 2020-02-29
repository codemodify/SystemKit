package tests

import (
	"fmt"
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
)

func Test_01(t *testing.T) {
	logging.Instance().LogTrace("Trace line")
	logging.Instance().LogPanic("Panic line")
	logging.Instance().LogFatal("Fatal line")
	logging.Instance().LogError("Error line")
	logging.Instance().LogWarning("Warning line")
	logging.Instance().LogInfo("Info line")
	logging.Instance().LogSuccess("Success line")
	logging.Instance().LogDebug("Debug line")

	fmt.Println()
}
