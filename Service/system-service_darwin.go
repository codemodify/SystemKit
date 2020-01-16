// +build darwin

package Service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	helpersExec "github.com/codemodify/SystemKit/Helpers"
	helpersFiles "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	helpersUser "github.com/codemodify/SystemKit/Helpers"
	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

// MacOSService - Represents Mac OS Service service
type MacOSService struct {
	command ServiceCommand
}

// New -
func New(command ServiceCommand) SystemService {
	// override some values - platform specific
	// https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html
	logDir := filepath.Join(helpersUser.HomeDir(""), "Library/Logs", command.Name)
	if helpersUser.IsRoot() {
		logDir = filepath.Join("/Library/Logs", command.Name)
	}

	command.Args = append([]string{command.Executable}, command.Args...)
	command.KeepAlive = true
	command.RunAtLoad = true
	command.StdOutPath = filepath.Join(logDir, command.Name+".stdout.log")
	command.StdErrPath = filepath.Join(logDir, command.Name+".stderr.log")

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
	dir := filepath.Dir(thisRef.FilePath())
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("making sure folder exists: ", dir),
	})
	os.MkdirAll(dir, os.ModePerm)

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("generating plist file"),
	})
	fileContent, err := thisRef.FileContent()
	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("writing plist to: ", thisRef.FilePath()),
	})
	err = ioutil.WriteFile(thisRef.FilePath(), fileContent, 0644)
	if err != nil {
		return err
	}

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
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("loading plist with launchctl"),
	})

	cmd := "launchctl"
	args := []string{"load", "-w", thisRef.FilePath()}
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("RUNNING: %s %s", cmd, strings.Join(args, " ")),
	})
	_, err := helpersExec.ExecWithArgs(cmd, args...)
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
	cmd := "launchctl"
	args := []string{"unload", "-w", thisRef.FilePath()}
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("RUNNING: %s %s", cmd, strings.Join(args, " ")),
	})
	_, err := helpersExec.ExecWithArgs(cmd, args...)
	if err != nil {
		e := strings.ToLower(err.Error())

		if strings.Contains(e, "could not find specified service") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("no service matching plist running: ", thisRef.command.DisplayLabel),
			})
			return nil
		}

		if strings.Contains(e, "no such file or directory") {
			logging.Instance().LogInfoWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprint("plist file doesn't exist, nothing to stop: ", thisRef.command.DisplayLabel),
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

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("remove plist file"),
	})
	err = os.Remove(thisRef.FilePath())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such file or directory") {
			return nil
		}

		return err
	}

	// sudo launchctl
	cmd := "launchctl"
	args := []string{"remove", thisRef.command.Name}
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("RUNNING: %s %s", cmd, strings.Join(args, " ")),
	})
	_, err = helpersExec.ExecWithArgs(cmd, args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error getting launchctl status: ", err),
		})
		return err
	}

	return nil
}

// Status -
func (thisRef MacOSService) Status() (ServiceStatus, error) {
	cmd := "launchctl"
	args := []string{"list"}
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("RUNNING: %s %s", cmd, strings.Join(args, " ")),
	})
	list, err := helpersExec.ExecWithArgs(cmd, args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error getting launchctl status: ", err),
		})
		return ServiceStatus{}, err
	}

	lines := strings.Split(strings.TrimSpace(string(list)), "\n")
	if thisRef.command.DisplayLabel == "" {
		return ServiceStatus{}, err
	}

	status := ServiceStatus{}

	for _, line := range lines {

		// logger.Log("line: ", line)

		chunks := strings.Split(line, "\t")

		if chunks[2] == thisRef.command.DisplayLabel {
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
	return helpersFiles.FileOrFolderExists(thisRef.FilePath())
}

// FilePath -
func (thisRef MacOSService) FilePath() string {
	if helpersUser.IsRoot() {
		return filepath.Join("/Library/LaunchDaemons", thisRef.command.Name+".plist")
	}

	return filepath.Join(helpersUser.HomeDir(""), "Library/LaunchAgents", thisRef.command.Name+".plist")
}

// FileContent -
func (thisRef MacOSService) FileContent() ([]byte, error) {
	plistTemplate := template.Must(template.New("launchdConfig").Parse(
		`<?xml version='1.0' encoding='UTF-8'?>
		<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
		<plist version='1.0'>
			<dict>
				<key>Label</key>
				<string>{{ .DisplayLabel }}</string>

				<key>ProgramArguments</key>
				<array>{{ range $arg := .Args }}
					<string>{{ $arg }}</string>{{ end }}
				</array>

				<key>StandardOutPath</key>
				<string>{{ .StdOutPath }}</string>

				<key>StandardErrorPath</key>
				<string>{{ .StdErrPath }}</string>

				<key>KeepAlive</key> <{{ .KeepAlive }}/>
				<key>RunAtLoad</key> <{{ .RunAtLoad }}/>

				<key>WorkingDirectory</key>
				<string>{{ .WorkingDirectory }}</string>
			</dict>
		</plist>
		`))

	var plistTemplateBytes bytes.Buffer
	if err := plistTemplate.Execute(&plistTemplateBytes, thisRef.command); err != nil {
		return nil, err
	}

	return plistTemplateBytes.Bytes(), nil
}
