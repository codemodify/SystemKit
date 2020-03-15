package Helpers

import (
	"os/exec"
	"runtime"
)

// Get the underlying OS command shell
func getOSC() string {

	osc := "sh"
	if runtime.GOOS == "windows" {
		osc = "cmd"
	}

	return osc
}

// Get the shell/command startup option to execute commands
func getOSE() string {

	ose := "-c"
	if runtime.GOOS == "windows" {
		ose = "/c"
	}
	return ose
}

// ExecutableExists -
func ExecutableExists(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Exec -
func Exec(command string) (string, error) {
	return ExecInFolder(command, "")
}

// ExecInFolder -
func ExecInFolder(command string, folder string) (string, error) {
	osc := getOSC()
	ose := getOSE()

	cmd := exec.Command(osc, ose, command)
	if !IsNullOrEmpty(folder) {
		cmd.Dir = folder
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ExecWithArgs -
func ExecWithArgs(name string, args ...string) (out string, err error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	return string(output), err
}
