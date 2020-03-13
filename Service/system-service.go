package Service

import (
	"encoding/json"
	"fmt"
)

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
	RunAsUser           string
	RunAsGroup          string
	OnStopDelegate      func()
}

func (thisRef ServiceCommand) String() string {
	bytes, err := json.Marshal(thisRef)
	if err != nil {
		// INFO: in normal app you could log this
		return ""
	}

	return string(bytes)
}

// ServiceStatus - is a generic representation of the service running on the system
type ServiceStatus struct {
	IsRunning bool
	PID       int
}

// ServiceErrorType -
type ServiceErrorType int

// ServiceErrorSuccess -
const (
	ServiceErrorSuccess      ServiceErrorType = iota
	ServiceErrorDoesNotExist                  = iota
	ServiceErrorOther                         = iota
)

func (thisRef ServiceErrorType) String() string {
	switch thisRef {
	case ServiceErrorSuccess:
		return "Success"

	case ServiceErrorDoesNotExist:
		return "Service Does Not Exist"

	case ServiceErrorOther:
		return "Other error occured"

	default:
		return fmt.Sprintf("%d", int(thisRef))
	}
}

// ServiceError -
type ServiceError struct {
	Type    ServiceErrorType
	Details error
}
