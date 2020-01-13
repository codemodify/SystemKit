// +build linux

package Service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	helpersFiles "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

// LinuxService - Represents Linux SystemD service
type LinuxService struct {
	command ServiceCommand
}

// New -
func New(command ServiceCommand) SystemService {
	return &LinuxService{
		command: command,
	}
}

// Run - is a no-op on Linux based systems
func (thisRef LinuxService) Run() error {
	return nil
}

// Install -
func (thisRef LinuxService) Install(start bool) error {
	unit := newSystemDFile(thisRef.command)

	path := unit.Path()
	dir := filepath.Dir(path)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("making sure folder exists: ", dir),
	})

	os.MkdirAll(dir, os.ModePerm)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("generating unit file"),
	})

	content, err := unit.Generate()

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("writing unit to: ", path),
	})

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("wrote unit:", content),
	})

	if start {
		err := thisRef.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

// Start -
func (thisRef LinuxService) Start() error {
	unit := newSystemDFile(thisRef.command)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "loading unit file with systemd",
	})

	_, err := runSystemCtlCommand("start", unit.Label)

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "enabling unit file with systemd",
	})

	_, err = runSystemCtlCommand("enable", unit.Label)

	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Created symlink") {
			return nil
		}
		return err
	}

	return nil
}

// Restart -
func (thisRef LinuxService) Restart() error {
	unit := newSystemDFile(thisRef.command)

	_, err := runSystemCtlCommand("reload-or-restart", unit.Label)

	if err != nil {
		return err
	}

	return nil
}

// Stop -
func (thisRef LinuxService) Stop() error {
	unit := newSystemDFile(thisRef.command)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "reloading daemon",
	})
	_, err := runSystemCtlCommand("daemon-reload", "")

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "stopping unit file with systemd",
	})

	_, err = runSystemCtlCommand("stop", unit.Label)
	// --force
	// --now

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "disabling unit file with systemd",
	})

	_, err = runSystemCtlCommand("disable", unit.Label)

	if err != nil {
		if strings.Contains(err.Error(), "Removed") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": "ignoring remove symlink error",
			})
			return nil
		}
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "reloading daemon",
	})

	_, err = runSystemCtlCommand("daemon-reload", "")

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "running reset-failed",
	})

	_, err = runSystemCtlCommand("reset-failed", "")

	if err != nil {
		return err
	}

	return nil
}

// Uninstall -
func (thisRef LinuxService) Uninstall() error {
	err := thisRef.Stop()

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "remove unit file",
	})

	unit := newSystemDFile(thisRef.command)
	err = unit.Remove()

	if err != nil {
		return err
	}

	return nil
}

// Status -
func (thisRef LinuxService) Status() (ServiceStatus, error) {
	unit := newSystemDFile(thisRef.command)
	active, _ := runSystemCtlCommand("is-active", unit.Label)

	status := ServiceStatus{}

	// Check if service is running
	if !strings.Contains(active, "active") {
		return status, nil
	}

	stat, _ := runSystemCtlCommand("status", unit.Label)

	// Get the PID from the status output
	lines := strings.Split(stat, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Main PID") {
			parts := strings.Split(strings.TrimSpace(line), " ")
			pid, _ := strconv.Atoi(parts[2])
			if pid != 0 {
				status.PID = pid
			}
		}
	}

	status.Running = true

	return status, nil
}

// Exists -
func (thisRef LinuxService) Exists() bool {
	unit := newSystemDFile(thisRef.command)
	return helpersFiles.FileOrFolderExists(unit.Path())
}
