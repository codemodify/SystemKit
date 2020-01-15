// +build darwin

package Service

import (
	"fmt"
	"strings"

	helpersExec "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

func runLaunchCtlCommand(args ...string) (out string, err error) {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("running command: launchctl ", strings.Join(args, " ")),
	})
	return helpersExec.ExecWithArgs("launchctl", args...)
}
