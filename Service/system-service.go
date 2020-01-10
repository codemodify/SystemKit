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
}

// ServiceCommand - what to execute as service
type ServiceCommand struct {
	Name          string
	DisplayLabel  string
	Executable    string
	WorkingDir    string
	Args          []string
	Description   string
	Documentation string // URL to your service documentation.
	Debug         bool   // Whether or not to turn on debug behavior
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
