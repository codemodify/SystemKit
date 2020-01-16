package Service

import "strings"

// SystemService - represents a generic system service configuration
type SystemService interface {
	Run() error
	Install(start bool) error
	Start() error
	Restart() error
	Stop() error
	Uninstall() error
	Status() (ServiceStatus, error)
	Exists() bool
	FilePath() string
	FileContent() ([]byte, error)
}

// ServiceCommand - What to execute as service
// These fields are a "common sense" mix of fields from SystemD and LaunchD.
// Some may be ignored on one or other platform but the implemetnation will
// try the max possible to respect the requested
type ServiceCommand struct {
	Name                string // usually this will be the file name
	DisplayLabel        string
	Description         string
	DocumentationURL    string
	Executable          string
	Args                []string
	WorkingDirectory    string
	Debug               bool
	KeepAlive           bool
	RunAtLoad           bool
	StdOutPath          string
	StdErrPath          string
	StartDelayInSeconds int
}

func (thisRef ServiceCommand) String() string {
	val := thisRef.Executable

	if len(thisRef.Args) > 0 {
		val = val + " " + strings.Join(thisRef.Args, " ")
	}

	return val
}

// ServiceStatus is a generic representation of the service running on the system
type ServiceStatus struct {
	Running bool
	PID     int
}
