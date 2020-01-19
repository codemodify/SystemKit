package List

type Process interface {
	PID() int
	ParentPID() int
	Executable() string
}

// Processes - returns a snapshot of all processes
func Processes() ([]Process, error) {
	return processes()
}

// FindProcessByPID - finds process by PID
func FindProcess(pid int) (Process, error) {
	return findProcess(pid)
}
