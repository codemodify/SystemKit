// +build darwin

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

// MacOSService - Represents Mac OS Service service
type MacOSService struct {
	command ServiceCommand
}

// New -
func New(command ServiceCommand) SystemService {
	return &MacOSService{
		command: command,
	}
}

// Run - is a no-op on Mac based systems
func (thisRef MacOSService) Run() error {
	return nil
}

// Install -
func (thisRef MacOSService) Install(start bool) error {
	plist := newPlist(thisRef.command)

	path := plist.Path()
	dir := filepath.Dir(path)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("making sure folder exists: ", dir),
	})

	os.MkdirAll(dir, os.ModePerm)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("generating plist file"),
	})

	content, err := plist.Generate()

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("writing plist to: ", path),
	})

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("wrote plist: ", content),
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
func (thisRef MacOSService) Start() error {
	plist := newPlist(thisRef.command)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("loading plist with launchctl"),
	})

	_, err := runLaunchCtlCommand("load", "-w", plist.Path())

	if err != nil {
		e := strings.ToLower(err.Error())

		// If not installed, install the service and then run start again.
		if strings.Contains(e, "no such file or directory") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("service not installed yet, installing..."),
			})

			err = thisRef.Install(true)

			if err != nil {
				return err
			}
		}

		// We don't care if the process fails because it is already
		// loaded
		if strings.Contains(e, "service already loaded") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("service already loaded"),
			})
			return nil
		}

		return err
	}

	return nil
}

// Restart -
func (thisRef MacOSService) Restart() error {
	err := thisRef.Stop()

	if err != nil {
		return err
	}

	err = thisRef.Start()

	if err != nil {
		return err
	}

	return nil
}

// Stop -
func (thisRef MacOSService) Stop() error {
	plist := newPlist(thisRef.command)

	_, err := runLaunchCtlCommand("unload", "-w", plist.Path())

	if err != nil {
		e := strings.ToLower(err.Error())

		if strings.Contains(e, "could not find specified service") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("no service matching plist running: ", plist.Label),
			})
			return nil
		}

		if strings.Contains(e, "no such file or directory") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("plist file doesn't exist, nothing to stop: ", plist.Label),
			})
			return nil
		}

		return err
	}

	return nil
}

// Uninstall -
func (thisRef MacOSService) Uninstall() error {
	err := thisRef.Stop()

	if err != nil {
		// If there is no matching process, don't throw an error
		// as it is already stopped.
		if strings.Contains(err.Error(), "exit status 3") != true {
			return err
		}
	}

	plist := newPlist(thisRef.command)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("remove plist file"),
	})

	err = os.Remove(plist.Path())

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such file or directory") {
			return nil
		}

		return err
	}

	return nil
}

// Status -
func (thisRef MacOSService) Status() (ServiceStatus, error) {
	plist := newPlist(thisRef.command)

	list, err := runLaunchCtlCommand("list")

	status := ServiceStatus{}

	if err != nil {
		logging.Instance().LogInfoWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error getting launchctl status: ", err),
		})
		return status, err
	}

	lines := strings.Split(strings.TrimSpace(string(list)), "\n")
	pattern := plist.Label

	if pattern == "" {
		return status, err
	}

	// logger.Log("running services:")

	for _, line := range lines {

		// logger.Log("line: ", line)

		chunks := strings.Split(line, "\t")

		if chunks[2] == pattern {
			if chunks[0] != "-" {
				pid, err := strconv.Atoi(chunks[0])

				if err != nil {
					return status, err
				}
				status.PID = pid
			}

			if status.PID != 0 {
				status.Running = true
			}
		}
	}

	return status, nil
}

// Exists -
func (thisRef MacOSService) Exists() bool {
	plist := newPlist(thisRef.command)

	return helpersFiles.FileOrFolderExists(plist.Path())
}
