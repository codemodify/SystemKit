package Monitor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"

	processList "github.com/codemodify/SystemKit/Processes/List"
	helpersJSON "github.com/codemodify/systemkit-helpers"
	helpersReflect "github.com/codemodify/systemkit-helpers"
	helpersStrings "github.com/codemodify/systemkit-helpers"
	logging "github.com/codemodify/systemkit-logging"
	loggingC "github.com/codemodify/systemkit-logging/contracts"
)

const procMonLogID = "PROC-MON"

// TheProcessMonitor - Represents Windows service
type TheProcessMonitor struct {
	procs     map[string]Process
	procsInfo map[string]processInfo
	procsSync *sync.RWMutex
}

type processInfo struct {
	osCmd     *exec.Cmd
	startedAt time.Time
	stoppedAt time.Time
	pid       int
	err       error
}

// New -
func New() ProcessMonitor {
	return &TheProcessMonitor{
		procs:     map[string]Process{},
		procsInfo: map[string]processInfo{},
		procsSync: &sync.RWMutex{},
	}
}

// Spawn -
func (thisRef *TheProcessMonitor) Spawn(id string, process Process) error {
	pi := processInfo{
		osCmd:     exec.Command(process.Executable, process.Args...),
		startedAt: time.Now(),
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: preparing to spawn [%s], details [%s]", procMonLogID, id, helpersJSON.AsJSONString(process)),
	})

	// set working folder
	if !helpersStrings.IsNullOrEmpty(process.WorkingDirectory) {
		pi.osCmd.Dir = process.WorkingDirectory
	}

	// set env
	if process.Env != nil {
		pi.osCmd.Env = process.Env
	}

	// set stderr and stdout
	stdOutPipe, err := pi.osCmd.StdoutPipe()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: failed to get StdoutPipe for [%s], details [%v]", procMonLogID, process.Executable, err),
		})

		return err
	}

	stdErrPipe, err := pi.osCmd.StderrPipe()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s, failed to get StderrPipe for [%s], details [%v]", procMonLogID, process.Executable, err),
		})

		return err
	}

	if process.OnStdOut != nil {
		go readStdOutFromProc(stdOutPipe, process)
	}
	if process.OnStdErr != nil {
		go readStdErrFromProc(stdErrPipe, process)
	}

	// MODIFY
	thisRef.procsSync.Lock()
	thisRef.procs[id] = process
	thisRef.procsInfo[id] = pi
	thisRef.procsSync.Unlock()

	return thisRef.Start(id)
}

// Start -
func (thisRef *TheProcessMonitor) Start(id string) error {
	if thisRef.GetProcessInfo(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procsInfo[id]; !ok {
		return fmt.Errorf("ID %s, CHECK-IF-EXISTS failed", id)
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to start [%s]", procMonLogID, id),
	})

	piToUpdate := thisRef.procsInfo[id]
	delete(thisRef.procsInfo, id)

	err := piToUpdate.osCmd.Start()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error starting [%v], details [%s]", procMonLogID, thisRef.procs[id], err),
		})

		piToUpdate.err = err
		piToUpdate.stoppedAt = time.Now()

		return err
	}

	// MODIFY
	piToUpdate.pid = piToUpdate.osCmd.Process.Pid
	thisRef.procsInfo[id] = piToUpdate

	return nil
}

// Stop -
func (thisRef *TheProcessMonitor) Stop(id string) error {
	if !thisRef.GetProcessInfo(id).IsRunning() {
		return nil
	}

	thisRef.procsSync.RLock()
	procInfo := thisRef.procsInfo[id]
	thisRef.procsSync.RUnlock()

	count := 0
	maxStopAttempts := 20
	for {
		count++
		if count > maxStopAttempts {
			procInfo.err = fmt.Errorf("%s: can't stop [%s]", procMonLogID, id)
			break
		}

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: attempt #%d to stop [%s]", procMonLogID, count, id),
		})

		if !thisRef.GetProcessInfo(id).IsRunning() {
			break
		}

		procInfo.osCmd.Process.Signal(syscall.SIGINT)
		procInfo.osCmd.Process.Signal(syscall.SIGTERM)
		procInfo.osCmd.Process.Signal(syscall.SIGKILL)
		processKillHelper(procInfo.osCmd.Process.Pid)

		err := procInfo.osCmd.Process.Kill()
		if err != nil {
			procInfo.err = err
		}

		time.Sleep(500 * time.Millisecond)
		procInfo.osCmd.Process.Wait()
	}

	procInfo.stoppedAt = time.Now()

	// MODIFY
	thisRef.procsSync.Lock()
	delete(thisRef.procsInfo, id)
	thisRef.procsInfo[id] = procInfo
	thisRef.procsSync.Unlock()

	return procInfo.err
}

// Restart -
func (thisRef TheProcessMonitor) Restart(id string) error {
	err := thisRef.Stop(id)
	if err != nil {
		return err
	}

	return thisRef.Start(id)
}

// StopAll -
func (thisRef TheProcessMonitor) StopAll() []error {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to stop all", procMonLogID),
	})

	allErrors := []error{}

	for k := range thisRef.procs {
		allErrors = append(allErrors, thisRef.Stop(k))
	}

	return allErrors
}

// GetProcessInfo -
func (thisRef TheProcessMonitor) GetProcessInfo(id string) ProcessInfo {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	// CHECK-IF-EXISTS
	if _, ok := thisRef.procsInfo[id]; !ok {
		return processInfo{}
	}

	return thisRef.procsInfo[id]
}

// RemoveFromMonitor -
func (thisRef *TheProcessMonitor) RemoveFromMonitor(id string) {
	thisRef.procsSync.Lock()
	defer thisRef.procsSync.Unlock()

	if _, ok := thisRef.procs[id]; ok {
		delete(thisRef.procs, id) // delete
	}

	if _, ok := thisRef.procsInfo[id]; ok {
		delete(thisRef.procsInfo, id) // delete
	}
}

// GetAll -
func (thisRef TheProcessMonitor) GetAll() []string {
	thisRef.procsSync.RLock()
	defer thisRef.procsSync.RUnlock()

	allIds := []string{}
	for k := range thisRef.procs {
		allIds = append(allIds, k)
	}

	return allIds
}

func readStdOutFromProc(readerCloser io.ReadCloser, process Process) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: starting to read StdOut [%s]", procMonLogID, process.Executable),
	})

	// output := make([]byte, 5000)

	reader := bufio.NewReader(readerCloser)
	// lengthRead, err := reader.Read(output)
	line, _, err := reader.ReadLine()
	for err != io.EOF {
		// process.OnStdOut(output[0:lengthRead])
		process.OnStdOut(line)
		// lengthRead, err = reader.Read(output)
		line, _, err = reader.ReadLine()
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error reading StdOut [%s], details [%s]", procMonLogID, process.Executable, err.Error()),
		})
	}
}

func readStdErrFromProc(readerCloser io.ReadCloser, process Process) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: starting to read StdErr [%s]", procMonLogID, process.Executable),
	})

	// output := make([]byte, 5000)

	reader := bufio.NewReader(readerCloser)
	// lengthRead, err := reader.Read(output)
	line, _, err := reader.ReadLine()
	for err != io.EOF {
		// process.OnStdOut(output[0:lengthRead])
		process.OnStdOut(line)
		// lengthRead, err = reader.Read(output)
		line, _, err = reader.ReadLine()
	}

	if err != nil {
		logging.Instance().LogWarningWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error reading StdErr [%s], details [%s]", procMonLogID, process.Executable, err.Error()),
		})
	}
}

func (thisRef processInfo) IsRunning() bool {
	if thisRef.osCmd == nil || thisRef.osCmd.Process == nil {
		return false
	}

	p, err := processList.FindProcess(thisRef.osCmd.Process.Pid)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error finding process [%d], details [%v]", procMonLogID, thisRef.osCmd.Process.Pid, err.Error()),
		})

		return false
	}

	return p != nil
}

func (thisRef processInfo) ExitCode() int {
	if thisRef.osCmd == nil || thisRef.osCmd.ProcessState == nil {
		return 0
	}

	return thisRef.osCmd.ProcessState.ExitCode()
}

func (thisRef processInfo) StartedAt() time.Time {
	if thisRef.osCmd == nil {
		return time.Unix(0, 0)
	}

	return thisRef.startedAt
}

func (thisRef processInfo) StoppedAt() time.Time {
	if thisRef.osCmd == nil {
		return time.Unix(0, 0)
	}

	return thisRef.stoppedAt
}

func (thisRef processInfo) PID() int {
	if thisRef.osCmd == nil {
		return 0
	}

	return thisRef.pid
}
