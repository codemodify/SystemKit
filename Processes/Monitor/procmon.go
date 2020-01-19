package Monitor

import "time"

import "encoding/json"

// ProcessMonitor - represents a generic system service configuration
type ProcessMonitor interface {
	Spawn(id string, command Process) error
	Start(id string) error
	Stop(id string) error
	Restart(id string) error
	StopAll() []error
	GetProcessInfo(id string) ProcessInfo
	RemoveFromMonitor(id string)
	GetAll() []string
}

// ProcessOutputReader -
type ProcessOutputReader func([]byte)

// Process -
type Process struct {
	Executable          string
	Args                []string
	WorkingDirectory    string
	Env                 []string
	DelayStartInSeconds int
	RestartRetryCount   int // -1 means unlimited
	OnStdOut            ProcessOutputReader
	OnStdErr            ProcessOutputReader
}

// String - stringer interface
func (thisRef Process) String() string {
	data, _ := json.Marshal(thisRef)
	return string(data)
}

// ProcessInfo -
type ProcessInfo interface {
	IsRunning() bool
	ExitCode() int
	StartedAt() time.Time
	StoppedAt() time.Time
}
