// +build linux

package Service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	helpersExec "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	helpersUser "github.com/codemodify/SystemKit/Helpers"
)

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

type systemDFile struct {
	Label         string
	Executable    string
	WorkingDir    string
	Description   string
	Documentation string
	StdOutPath    string
	StdErrPath    string
	User          string
}

func newSystemDFile(command ServiceCommand) systemDFile {
	label := command.Name

	user := helpersUser.UserName("")
	if helpersUser.IsRoot() {
		user = "root"
	}

	unit := systemDFile{
		Label:         label,
		Executable:    command.String(),
		WorkingDir:    command.WorkingDir,
		Description:   command.Description,
		Documentation: command.Documentation,
		User:          user,
	}

	return unit
}

func (thisRef systemDFile) Generate() (string, error) {
	var tmpl bytes.Buffer
	t := template.Must(template.New("systemDFile").Parse(unitFileTemplate()))
	if err := t.Execute(&tmpl, thisRef); err != nil {
		return "", err
	}

	return tmpl.String(), nil
}

func (thisRef systemDFile) Path() string {
	file := thisRef.Label + ".service"

	if helpersUser.IsRoot() {
		return filepath.Join("/etc/systemd/system", file)
	}

	return filepath.Join(helpersUser.HomeDir(""), ".config/systemd/user", file)
}

func (thisRef systemDFile) Remove() error {
	return os.Remove(thisRef.Path())
}

func unitFileTemplate() string {
	return `[Unit]
After=network.target
Description={{ .Description }}
Documentation={{ .Documentation }}

[Service]
ExecStart={{ .Executable }}
WorkingDirectory={{ .WorkingDir }}
Restart=on-failure
Type=simple

[Install]
WantedBy=multi-user.target
`
}
