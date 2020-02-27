package Monitor

import (
	"time"
)

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
	RestartRetryCount   int                 // -1 means unlimited
	OnStdOut            ProcessOutputReader `json:"-"`
	OnStdErr            ProcessOutputReader `json:"-"`
}

// ProcessInfo -
type ProcessInfo interface {
	IsRunning() bool
	ExitCode() int
	StartedAt() time.Time
	StoppedAt() time.Time
	PID() int
}
