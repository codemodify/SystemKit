// +build windows

package Service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	helpersExec "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

var logTag = "SYSTEM-SERVICE"

// WindowsService - Represents Windows service
type WindowsService struct {
	command ServiceCommand
}

// New -
func New(command ServiceCommand) SystemService {
	return &WindowsService{
		command: command,
	}
}

// Run -
func (thisRef *WindowsService) Run() error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to run: %s", logTag, thisRef.command.Name),
	})

	wg := sync.WaitGroup{}

	wg.Add(1)
	var err error
	go func() {
		err = svc.Run(thisRef.command.Name, thisRef)
		wg.Done()
	}()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: running: %s", logTag, thisRef.command.Name),
	})
	wg.Wait()

	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: failed to run: %s, %v", logTag, thisRef.command.Name, err),
		})
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: stopped: %s", logTag, thisRef.command.Name),
	})

	return nil
}

// Install -
func (thisRef *WindowsService) Install(start bool) error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to install: %s", logTag, thisRef.command.Name),
	})

	// check if service exists
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: check if exists: %s", logTag, thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err.Type == ServiceErrorDoesNotExist { // this is a good thing
		if winService != nil {
			winService.Close()
		}
		if winServiceManager != nil {
			winServiceManager.Disconnect()
		}

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: does not exist: %s", logTag, thisRef.command.Name),
		})
	} else {
		if winService != nil {
			winService.Close()
		}
		if winServiceManager != nil {
			winServiceManager.Disconnect()
		}

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: service '%s' already exists: ", logTag, thisRef.command.Name),
		})

		return fmt.Errorf("service '%s' already exists: ", thisRef.command.Name)
	}

	// connect again as the `winServiceManager, winService` are null from previous step
	winServiceManager, winService, err = connectAndOpenService(thisRef.command.Name)
	if err.Type != ServiceErrorDoesNotExist {
		winServiceManager.Disconnect()

		return err.Details
	}
	defer winServiceManager.Disconnect()

	// Create the system service
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: creating: '%s', binary: '%s', args: '%s'", logTag, thisRef.command.Name, thisRef.command.Executable, thisRef.command.Args),
	})

	if winService != nil {
		winService.Close()
	}

	winService, err1 := winServiceManager.CreateService(
		thisRef.command.Name,
		thisRef.command.Executable,
		mgr.Config{
			StartType:   mgr.StartAutomatic,
			DisplayName: thisRef.command.Name,
			Description: thisRef.command.Description,
		},
		thisRef.command.Args...,
	)
	if err1 != nil {
		if winService != nil {
			winService.Close()
		}
		if winServiceManager != nil {
			winServiceManager.Disconnect()
		}

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error creating: %s, details: %v ", logTag, thisRef.command.Name, err1),
		})

		return err1
	}

	winService.Close()
	winServiceManager.Disconnect()

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: created: '%s', binary: '%s', args: '%s'", logTag, thisRef.command.Name, thisRef.command.Executable, thisRef.command.Args),
	})

	if start {
		if err := thisRef.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Start -
func (thisRef *WindowsService) Start() error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to start: ", logTag, thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err.Type != ServiceErrorSuccess {
		return err.Details
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	err1 := winService.Start() // thisRef.command.Args...
	if err1 != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error starting: %s, %v", logTag, thisRef.command.Name, err1),
		})

		return fmt.Errorf("error starting: %s, %v", thisRef.command.Name, err1)
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: started: %s", logTag, thisRef.command.Name),
	})

	return nil
}

// Restart -
func (thisRef *WindowsService) Restart() error {
	if err := thisRef.Stop(); err != nil {
		return err
	}

	if err := thisRef.Start(); err != nil {
		return err
	}

	return nil
}

// Stop -
func (thisRef *WindowsService) Stop() error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to stop: %s", logTag, thisRef.command.Name),
	})

	if thisRef.command.OnStopDelegate != nil {
		logger.Debug("ACTUAL: OnStopDelegate() beferring-calling")

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: OnStopDelegate beferring-calling: %s", logTag, thisRef.command.Name),
		})

		thisRef.command.OnStopDelegate()

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: OnStopDelegate after-calling: %s", logTag, thisRef.command.Name),
		})
	}

	err := thisRef.control(svc.Stop, svc.Stopped)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "service does not exist") {
			return nil
		} else if strings.Contains(e, "service has not been started") {
			return nil
		}

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error %s, details: %s", logTag, thisRef.command.Name, err.Error()),
		})

		return err
	}

	attempt := 0
	maxAttempts := 10
	wait := 3 * time.Second
	for {
		attempt++

		logging.Instance().LogDebugWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: waiting for service to stop", logTag),
		})

		// Wait a few seconds before retrying
		time.Sleep(wait)

		// Attempt to start the service again
		stat, err := thisRef.Status()
		if err != nil {
			return err
		}

		// If it is now running, exit the retry loop
		if !stat.IsRunning {
			break
		}

		if attempt == maxAttempts {
			return errors.New("could not stop system service after multiple attempts")
		}
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: stopped: %s", logTag, thisRef.command.Name),
	})

	return nil
}

// Uninstall -
func (thisRef *WindowsService) Uninstall() error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to uninstall: %s", logTag, thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err.Type != ServiceErrorSuccess {
		return err.Details
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	err1 := winService.Delete()
	if err1 != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: failed to uninstall: %s, %v", logTag, thisRef.command.Name, err1),
		})

		return err1
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: uninstalled: %s", logTag, thisRef.command.Name),
	})

	return nil
}

// Status -
func (thisRef *WindowsService) Status() (ServiceStatus, error) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: querying status: %s", logTag, thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err.Type != ServiceErrorSuccess {
		return ServiceStatus{}, err.Details
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	stat, err1 := winService.Query()
	if err1 != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("%s: error getting service status: ", logTag, err1),
		})

		return ServiceStatus{}, fmt.Errorf("error getting service status: %v", err1)
	}

	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: service status: %#v", logTag, stat),
	})

	return ServiceStatus{
		PID:       int(stat.ProcessId),
		IsRunning: stat.State == svc.Running,
	}, nil
}

// Exists -
func (thisRef *WindowsService) Exists() bool {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: checking existance: %s", logTag, thisRef.command.Name),
	})

	args := []string{"queryex", fmt.Sprintf("\"%s\"", thisRef.command.Name)}

	// https://www.computerhope.com/sc-command.htm
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: running: 'sc %s'", logTag, strings.Join(args, " ")),
	})

	_, err := helpersExec.ExecWithArgs("sc", args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error when checking %s: ", logTag, err),
		})

		return false
	}

	return true
}

// FilePath -
func (thisRef *WindowsService) FilePath() string {
	return ""
}

// FileContent -
func (thisRef *WindowsService) FileContent() ([]byte, error) {
	return []byte{}, nil
}

// Execute - implement the Windows `service.Handler` interface
func (thisRef *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: WINDOWS SERVICE EXECUTE", logTag),
	})

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus

			case svc.Stop, svc.Shutdown:
				if thisRef.command.OnStopDelegate != nil {
					logging.Instance().LogDebugWithFields(loggingC.Fields{
						"method":  helpersReflect.GetThisFuncName(),
						"message": fmt.Sprintf("%s: OnStopDelegate beferring-calling: %s", logTag, thisRef.command.Name),
					})

					go thisRef.command.OnStopDelegate()

					logging.Instance().LogDebugWithFields(loggingC.Fields{
						"method":  helpersReflect.GetThisFuncName(),
						"message": fmt.Sprintf("%s: OnStopDelegate after-calling: %s", logTag, thisRef.command.Name),
					})
				}

				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				// testOutput := strings.Join(args, "-")
				// testOutput += fmt.Sprintf("-%d", c.Context)
				// logging.Instance().LogDebugWithFields(loggingC.Fields{
				// 	"method":  helpersReflect.GetThisFuncName(),
				// 	"message": fmt.Sprintf("%s: %", logTag, testOutput),
				// })

				break loop

			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}

			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

			default:
				logging.Instance().LogWarningWithFields(loggingC.Fields{
					"method":  helpersReflect.GetThisFuncName(),
					"message": fmt.Sprintf("%s: unexpected control request #%d", logTag, c),
				})
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func (thisRef *WindowsService) control(command svc.Cmd, state svc.State) error {
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: attempting to control: %s, cmd: %v", logTag, thisRef.command.Name, command),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err.Type != ServiceErrorSuccess {
		return err.Details
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	status, err1 := winService.Control(command)
	if err1 != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: could not send control: %d, to: %s, details: %v", logTag, command, thisRef.command.Name, err1),
		})

		return fmt.Errorf("could not send control: %d, to: %s, details: %v", command, thisRef.command.Name, err1)
	}

	timeout := time.Now().Add(10 * time.Second)
	for status.State != state {
		// Exit if a timeout is reached
		if timeout.Before(time.Now()) {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("%s: timeout waiting for service to go to state=%d", logTag, state),
			})

			return fmt.Errorf("timeout waiting for service to go to state=%d", state)
		}

		time.Sleep(300 * time.Millisecond)

		// Make sure transition happens to the desired state
		status, err1 = winService.Query()
		if err1 != nil {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("%s: could not retrieve service status: %v", logTag, err1),
			})

			return fmt.Errorf("could not retrieve service status: %v", err1)
		}
	}

	return nil
}

func connectAndOpenService(serviceName string) (*mgr.Mgr, *mgr.Service, ServiceError) {
	// connect to Windows Service Manager
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: connecting to Windows Service Manager", logTag),
	})

	winServiceManager, err := mgr.Connect()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error connecting to Windows Service Manager: %v", logTag, err),
		})
		return nil, nil, ServiceError{Type: ServiceErrorOther, Details: err}
	}

	// open service to manage it
	logging.Instance().LogDebugWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s: opening service: %s", logTag, serviceName),
	})

	winService, err := winServiceManager.OpenService(serviceName)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s: error opening service: %s, %v", logTag, serviceName, err),
		})

		return winServiceManager, nil, ServiceError{Type: ServiceErrorDoesNotExist, Details: err}
	}

	return winServiceManager, winService, ServiceError{Type: ServiceErrorSuccess}
}
