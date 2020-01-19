package Monitor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	helpersStrings "github.com/codemodify/SystemKit/Helpers"
	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
	processList "github.com/codemodify/SystemKit/Processes/List"
)

// WindowsProcessMonitor - Represents Windows service
type WindowsProcessMonitor struct {
	procs     map[string]Process
	procsInfo map[string]*processInfo
	procsSync sync.RWMutex
}

type processInfo struct {
	osCmd     *exec.Cmd
	startedAt time.Time
	stoppedAt time.Time
	err       error
}

// New -
func New() ProcessMonitor {
	return &WindowsProcessMonitor{
		procs:     map[string]Process{},
		procsInfo: map[string]*processInfo{},
		procsSync: sync.RWMutex{},
	}
}

// Spawn -
func (thisRef *WindowsProcessMonitor) Spawn(id string, process Process) error {
	thisRef.procsSync.Lock()

	thisRef.procs[id] = process
	thisRef.procsInfo[id] = &processInfo{
		osCmd:     exec.Command(thisRef.procs[id].Executable, thisRef.procs[id].Args...),
		startedAt: time.Now(),
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to spawn: %s", id),
	})
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to spawn details: %s", thisRef.procs[id]),
	})

	// set working folder
	if !helpersStrings.IsNullOrEmpty(thisRef.procs[id].WorkingDirectory) {
		thisRef.procsInfo[id].osCmd.Dir = thisRef.procs[id].WorkingDirectory
	}

	// set env
	if thisRef.procs[id].Env != nil {
		thisRef.procsInfo[id].osCmd.Env = thisRef.procs[id].Env
	}

	// set stderr and stdout
	stdOutPipe, err := thisRef.procsInfo[id].osCmd.StdoutPipe()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("failed to get StdoutPipe: %s, %v", thisRef.procs[id].Executable, err),
		})

		return err
	}

	stdErrPipe, err := thisRef.procsInfo[id].osCmd.StderrPipe()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("failed to get StderrPipe: %s, %v", thisRef.procs[id].Executable, err),
		})

		return err
	}

	if thisRef.procs[id].OnStdOut != nil {
		go readStdOutFromProc(stdOutPipe, thisRef.procs[id])
	}
	if thisRef.procs[id].OnStdErr != nil {
		go readStdErrFromProc(stdErrPipe, thisRef.procs[id])
	}

	thisRef.procsSync.Unlock()

	return thisRef.Start(id)
}

// Start -
func (thisRef *WindowsProcessMonitor) Start(id string) error {
	if thisRef.GetProcessInfo(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to start: %s", id),
	})

	err := thisRef.procsInfo[id].osCmd.Start()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("error starting: %v, details: %s", thisRef.procs[id], err),
		})

		thisRef.procsInfo[id].err = err
		thisRef.procsInfo[id].stoppedAt = time.Now()

		return err
	}

	return nil
}

// Stop -
func (thisRef *WindowsProcessMonitor) Stop(id string) error {
	if !thisRef.GetProcessInfo(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.RLock()
	procInfo := thisRef.procsInfo[id]
	thisRef.procsSync.RUnlock()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to stop: %s", id),
	})

	count := 0
	for {
		count++

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("attempt #%d", count),
		})

		if !thisRef.GetProcessInfo(id).IsRunning() {
			break
		}

		err := procInfo.osCmd.Process.Kill()
		if err != nil {
			procInfo.err = err
		}

		// err = syscall.Kill(procInfo.osCmd.Process.Pid, syscall.SIGKILL)
		// if err != nil {
		// 	procInfo.err = err
		// }

		time.Sleep(500 * time.Millisecond)
		procInfo.osCmd.Process.Wait()
	}

	procInfo.stoppedAt = time.Now()

	return procInfo.err
}

// Restart -
func (thisRef *WindowsProcessMonitor) Restart(id string) error {
	err := thisRef.Stop(id)
	if err != nil {
		return err
	}

	return thisRef.Start(id)
}

// StopAll -
func (thisRef *WindowsProcessMonitor) StopAll() []error {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to stop all"),
	})

	allErrors := []error{}

	for k := range thisRef.procs {
		allErrors = append(allErrors, thisRef.Stop(k))
	}

	return allErrors
}

// GetProcessInfo -
func (thisRef *WindowsProcessMonitor) GetProcessInfo(id string) ProcessInfo {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	return thisRef.procsInfo[id]
}

// RemoveFromMonitor -
func (thisRef *WindowsProcessMonitor) RemoveFromMonitor(id string) {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	if _, ok := thisRef.procs[id]; ok {
		delete(thisRef.procs, id) // delete
	}

	if _, ok := thisRef.procsInfo[id]; ok {
		delete(thisRef.procsInfo, id) // delete
	}
}

func readStdOutFromProc(readerCloser io.ReadCloser, process Process) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("starting to read StdOut: %s", process.Executable),
	})

	output := make([]byte, 5000)

	reader := bufio.NewReader(readerCloser)
	lengthRead, err := reader.Read(output)
	for err != io.EOF {
		process.OnStdOut(output[0:lengthRead])
		lengthRead, err = reader.Read(output)
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s", err),
		})
	}
}

func readStdErrFromProc(readerCloser io.ReadCloser, process Process) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("starting to read StdErr: %s", process.Executable),
	})

	output := make([]byte, 5000)

	reader := bufio.NewReader(readerCloser)
	lengthRead, err := reader.Read(output)
	for err != io.EOF {
		process.OnStdErr(output[0:lengthRead])
		lengthRead, err = reader.Read(output)
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s", err),
		})
	}
}

func (thisRef processInfo) IsRunning() bool {
	if thisRef.osCmd.Process == nil {
		return false
	}

	p, err := processList.FindProcess(thisRef.osCmd.Process.Pid)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%v", err),
		})

		return false
	}

	return p != nil
}

func (thisRef processInfo) ExitCode() int {
	if thisRef.osCmd.ProcessState == nil {
		return 0
	}

	return thisRef.osCmd.ProcessState.ExitCode()
}

func (thisRef processInfo) StartedAt() time.Time {
	return thisRef.startedAt
}

func (thisRef processInfo) StoppedAt() time.Time {
	return thisRef.stoppedAt
}
