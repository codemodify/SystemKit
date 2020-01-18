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
func (thisRef WindowsService) Run() error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to run: %s", thisRef.command.Name),
	})

	wg := sync.WaitGroup{}

	wg.Add(1)
	var err error
	go func() {
		err = svc.Run(thisRef.command.Name, &thisRef)
		wg.Done()
	}()

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("running: %s", thisRef.command.Name),
	})
	wg.Wait()

	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("failed to run: %s, %v", thisRef.command.Name, err),
		})
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("stopped: %s", thisRef.command.Name),
	})

	return nil
}

// Install -
func (thisRef WindowsService) Install(start bool) error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to install: %s", thisRef.command.Name),
	})

	// check if service exists
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("check if service exists: %s", thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err == nil {
		if winService != nil {
			winService.Close()
		}
		if winServiceManager != nil {
			winServiceManager.Disconnect()
		}

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("service '%s' already exists: ", thisRef.command.Name),
		})

		return fmt.Errorf("service '%s' already exists: ", thisRef.command.Name)
	}

	// Create the system service
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("creating service: '%s', binary: '%s', args: '%s'", thisRef.command.Name, thisRef.command.Executable, thisRef.command.Args),
	})

	conf := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: thisRef.command.Name,
		Description: thisRef.command.Description,
	}

	if winService != nil {
		winService.Close()
	}
	winService, err = winServiceManager.CreateService(thisRef.command.Name, thisRef.command.Executable, conf, thisRef.command.Args...)
	if err != nil {
		if winService != nil {
			winService.Close()
		}
		if winServiceManager != nil {
			winServiceManager.Disconnect()
		}

		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error creating service: ", err),
		})

		return err
	}

	winService.Close()
	winServiceManager.Disconnect()

	if start {
		if err := thisRef.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Start -
func (thisRef WindowsService) Start() error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("attempting to start: ", thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err != nil {
		return err
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	err = winService.Start(thisRef.command.Args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("error starting: %s, %v", thisRef.command.Name, err),
		})

		return fmt.Errorf("error starting: %s, %v", thisRef.command.Name, err)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("started: %s", thisRef.command.Name),
	})

	return nil
}

// Restart -
func (thisRef WindowsService) Restart() error {
	if err := thisRef.Stop(); err != nil {
		return err
	}

	if err := thisRef.Start(); err != nil {
		return err
	}

	return nil
}

// Stop -
func (thisRef WindowsService) Stop() error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("attempting to stop: ", thisRef.command.Name),
	})

	err := thisRef.control(svc.Stop, svc.Stopped)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "service does not exist") {
			return nil
		}

		return err
	}

	attempt := 0
	maxAttempts := 10
	wait := 3 * time.Second
	for {
		attempt++

		logging.Instance().LogInfoWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("waiting for service to stop"),
		})

		// Wait a few seconds before retrying
		time.Sleep(wait)

		// Attempt to start the service again
		stat, err := thisRef.Status()
		if err != nil {
			return err
		}

		// If it is now running, exit the retry loop
		if !stat.Running {
			break
		}

		if attempt == maxAttempts {
			return errors.New("could not stop system service after multiple attempts")
		}
	}

	return nil
}

// Uninstall -
func (thisRef WindowsService) Uninstall() error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to uninstall: %s", thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err != nil {
		return err
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	err = winService.Delete()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("failed to uninstall: %s, %v", thisRef.command.Name, err),
		})

		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("uninstalled: %s", thisRef.command.Name),
	})

	return nil
}

// Status -
func (thisRef WindowsService) Status() (ServiceStatus, error) {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("querying status: %s", thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err != nil {
		return ServiceStatus{}, err
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	stat, err := winService.Query()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error getting service status: ", err),
		})

		return ServiceStatus{}, fmt.Errorf("error getting service status: ", err)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("service status: %#v", stat),
	})

	return ServiceStatus{
		PID:     int(stat.ProcessId),
		Running: stat.State == svc.Running,
	}, nil
}

// Exists -
func (thisRef WindowsService) Exists() bool {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("checking existance: %s", thisRef.command.Name),
	})

	args := []string{"queryex", fmt.Sprintf("\"%s\"", thisRef.command.Name)}

	// https://www.computerhope.com/sc-command.htm
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("running: 'sc %s'", strings.Join(args, " ")),
	})

	_, err := helpersExec.ExecWithArgs("sc", args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("error when checking %s: ", err),
		})

		return false
	}

	return true
}

// FilePath -
func (thisRef WindowsService) FilePath() string {
	return ""
}

// FileContent -
func (thisRef WindowsService) FileContent() ([]byte, error) {
	return []byte{}, nil
}

// Execute - implement the Windows `service.Handler` interface
func (thisRef WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("WINDOWS SERVICE EXECUTE"),
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
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d", c.Context)
				logging.Instance().LogInfoWithFields(loggingC.Fields{
					"method":  helpersReflect.GetThisFuncName(),
					"message": fmt.Sprint(testOutput),
				})

				break loop

			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}

			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

			default:
				logging.Instance().LogWarningWithFields(loggingC.Fields{
					"method":  helpersReflect.GetThisFuncName(),
					"message": fmt.Sprintf("unexpected control request #%d", c),
				})
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func (thisRef WindowsService) control(command svc.Cmd, state svc.State) error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("attempting to control: %s", thisRef.command.Name),
	})

	winServiceManager, winService, err := connectAndOpenService(thisRef.command.Name)
	if err != nil {
		return err
	}
	defer winServiceManager.Disconnect()
	defer winService.Close()

	status, err := winService.Control(command)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("could not send control: %d, to: %s, details: %v", command, thisRef.command.Name, err),
		})

		return fmt.Errorf("could not send control: %d, to: %s, details: %v", command, thisRef.command.Name, err)
	}

	timeout := time.Now().Add(10 * time.Second)
	for status.State != state {
		// Exit if a timeout is reached
		if timeout.Before(time.Now()) {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("timeout waiting for service to go to state=%d", state),
			})

			return fmt.Errorf("timeout waiting for service to go to state=%d", state)
		}

		time.Sleep(300 * time.Millisecond)

		// Make sure transition happens to the desired state
		status, err = winService.Query()
		if err != nil {
			logging.Instance().LogErrorWithFields(loggingC.Fields{
				"method":  helpersReflect.GetThisFuncName(),
				"message": fmt.Sprintf("could not retrieve service status: %v", err),
			})

			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}

	return nil
}

func connectAndOpenService(serviceName string) (*mgr.Mgr, *mgr.Service, error) {
	// connect to Windows Service Manager
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("connecting to Windows Service Manager"),
	})

	winServiceManager, err := mgr.Connect()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("error connecting to Windows Service Manager: %v", err),
		})
		return nil, nil, err
	}

	// open service to manage it
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("opening service: %s", serviceName),
	})

	winService, err := winServiceManager.OpenService(serviceName)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("error opening service: %s, %v", serviceName, err),
		})

		return nil, nil, err
	}

	return winServiceManager, winService, nil
}
