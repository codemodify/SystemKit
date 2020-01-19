// +build linux

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
	path := thisRef.FilePath()
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

	content, err := thisRef.FileContent()

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("writing unit to: ", path),
	})

	err = ioutil.WriteFile(path, content, 0644)

	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("wrote unit: %s", string(content)),
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
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "loading unit file with systemd",
	})

	_, err := runSystemCtlCommand("start", thisRef.command.Name)
	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "enabling unit file with systemd",
	})

	_, err = runSystemCtlCommand("enable", thisRef.command.Name)
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
	_, err := runSystemCtlCommand("reload-or-restart", thisRef.command.Name)
	if err != nil {
		return err
	}

	return nil
}

// Stop -
func (thisRef LinuxService) Stop() error {
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
	_, err = runSystemCtlCommand("stop", thisRef.command.Name)
	if err != nil {
		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": "disabling unit file with systemd",
	})
	_, err = runSystemCtlCommand("disable", thisRef.command.Name)
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
	err = os.Remove(thisRef.FilePath())
	if err != nil {
		return err
	}

	return nil
}

// Status -
func (thisRef LinuxService) Status() (ServiceStatus, error) {
	active, _ := runSystemCtlCommand("is-active", thisRef.command.Name)

	status := ServiceStatus{}

	// Check if service is running
	if !strings.Contains(active, "active") {
		return status, nil
	}

	stat, _ := runSystemCtlCommand("status", thisRef.command.Name)

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
	return helpersFiles.FileOrFolderExists(thisRef.FilePath())
}

// FilePath -
func (thisRef LinuxService) FilePath() string {
	if helpersUser.IsRoot() {
		return filepath.Join("/etc/systemd/system", thisRef.command.Name+".service")
	}

	return filepath.Join(helpersUser.HomeDir(""), ".config/systemd/user", thisRef.command.Name+".service")
}

// FileContent -
func (thisRef LinuxService) FileContent() ([]byte, error) {
	transformedCommand := transformCommandForSaveToDisk(thisRef.command)

	systemDServiceFileTemplate := template.Must(template.New("systemDFile").Parse(
		`[Unit]
		After=network.target
		Description={{ .Description }}
		Documentation={{ .DocumentationURL }}
		
		[Service]
		ExecStart={{ .Executable }}
		WorkingDirectory={{ .WorkingDirectory }}
		Restart=on-failure
		Type=simple

		{{ if .RunAsUser }}
		User={{ .RunAsUser }}
		{{ end }}
		{{ if .RunAsGroup }}
		Group={{ .RunAsGroup }}
		{{ end }}
		
		[Install]
		WantedBy=multi-user.target
	`))

	var systemDServiceFileTemplateAsBytes bytes.Buffer
	if err := systemDServiceFileTemplate.Execute(&systemDServiceFileTemplateAsBytes, transformedCommand); err != nil {
		return nil, err
	}

	return systemDServiceFileTemplateAsBytes.Bytes(), nil
}

func runSystemCtlCommand(cmd string, label string) (out string, err error) {
	args := strings.Split(cmd, " ")

	if !helpersUser.IsRoot() {
		args = append(args, "--user")
	}

	if label != "" {
		args = append(args, label)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("running command: systemctl ", strings.Join(args, " ")),
	})

	return helpersExec.ExecWithArgs("systemctl", args...)
}

func transformCommandForSaveToDisk(command ServiceCommand) ServiceCommand {
	if len(command.Args) > 0 {
		command.Executable = fmt.Sprintf("%s %s", command.Executable, strings.Join(command.Args, " "))
	}

	return command
}
