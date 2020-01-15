// +build darwin

package Service

import (
	"bytes"
	"path/filepath"
	"text/template"

	helpersUser "github.com/codemodify/SystemKit/Helpers"
)

type pListFile struct {
	Label            string
	Program          string
	ProgramArguments []string
	KeepAlive        bool
	RunAtLoad        bool
	StdOutPath       string
	StdErrPath       string
}

func newPlist(command ServiceCommand) pListFile {
	logDir := filepath.Join(helpersUser.HomeDir(""), "Library/Logs", command.Name)
	if helpersUser.IsRoot() {
		logDir = filepath.Join("/Library/Logs", command.Name)
	}
	args := []string{command.Executable}
	if len(command.Args) != 0 {
		args = append(args, command.Args...)
	}

	pl := pListFile{
		Label:            command.DisplayLabel,
		ProgramArguments: args,
		KeepAlive:        true,
		RunAtLoad:        true,
		StdOutPath:       filepath.Join(logDir, command.Name+".stdout.log"),
		StdErrPath:       filepath.Join(logDir, command.Name+".stderr.log"),
	}

	return pl
}

// Generate -
func (thisRef pListFile) Generate() (string, error) {
	var tmpl bytes.Buffer
	t := template.Must(template.New("launchdConfig").Parse(plistTemplate()))
	if err := t.Execute(&tmpl, thisRef); err != nil {
		return "", err
	}

	return tmpl.String(), nil
}

// Path -
func (thisRef pListFile) Path() string {
	label := thisRef.Label + ".pListFile"
	if helpersUser.IsRoot() {
		return filepath.Join("/Library/LaunchDaemons/", label)
	}

	return filepath.Join(helpersUser.HomeDir(""), "Library/LaunchAgents/", label)
}

// plistTemplate - generates the contents of the pListFile file
func plistTemplate() string {
	return `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE pListFile PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<pListFile version='1.0'>
  <dict>
    <key>Label</key><string>{{ .Label }}</string>{{ if .Program }}
    <key>Program</key><string>{{ .Program }}</string>{{ end }}
    {{ if .ProgramArguments }}<key>ProgramArguments</key>
    <array>{{ range $arg := .ProgramArguments }}
      <string>{{ $arg }}</string>{{ end }}
    </array>{{ end }}
    <key>StandardOutPath</key>
    <string>{{ .StdOutPath }}</string>
    <key>StandardErrorPath</key>
    <string>{{ .StdErrPath }}</string>
    <key>KeepAlive</key> <{{ .KeepAlive }}/>
    <key>RunAtLoad</key> <{{ .RunAtLoad }}/>
  </dict>
</pListFile>
`
}
