// +build darwin

package Service

import "strings"

func runLaunchCtlCommand(args ...string) (out string, err error) {
	logger.Log("running command: launchctl ", strings.Join(args, " "))
	return runCommand("launchctl", args...)
}
