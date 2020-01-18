// +build windows

package Service

import (
	"errors"
	"fmt"
	"strings"
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
		"message": fmt.Sprintf("starting %s service", thisRef.command.Name),
	})

	err := svc.Run(thisRef.command.Name, &thisRef)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprintf("%s service failed: %v", thisRef.command.Name, err),
		})

		return err
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("%s service stopped", thisRef.command.Name),
	})

	return nil
}

// Install -
func (thisRef WindowsService) Install(start bool) error {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("installing system service: ", thisRef.command.Name),
	})

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error connecting to service manager: ", err),
		})
		return err
	}
	defer m.Disconnect()

	// Open the service so we can manage it
	srv, err := m.OpenService(thisRef.command.Name)
	if err == nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error opening the service: ", thisRef.command.Name),
		})
		srv.Close()
		return fmt.Errorf("service %s already exists", thisRef.command.Name)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("creating service \"%s\" at path \"%s\" with args \"%s\"", thisRef.command.Name, thisRef.command.Executable, thisRef.command.Args),
	})

	// Create the system service
	conf := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: thisRef.command.Name,
		Description: thisRef.command.Description,
	}
	srv, err = m.CreateService(thisRef.command.Name, thisRef.command.Executable, conf, thisRef.command.Args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error creating service: ", err),
		})
		return err
	}
	defer srv.Close()

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
		"message": fmt.Sprint("starting system service: ", thisRef.command.Name),
	})

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error connecting to service manager: ", err),
		})
		return fmt.Errorf("could not connect to service manager: %v", err)
	}
	defer m.Disconnect()

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("opening system service"),
	})

	// Open the service so we can manage it
	srv, err := m.OpenService(thisRef.command.Name)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error opening service: ", err),
		})
		return fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("attempting to start system service"),
	})

	err = srv.Start(thisRef.command.Args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error starting service: ", err),
		})
		return fmt.Errorf("could not start service: %v", err)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("running service"),
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
	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// Open the service so we can manage it
	srv, err := m.OpenService(thisRef.command.Name)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "not installed") || strings.Contains(e, "does not exist") {
			return nil
		}
		return err
	}
	defer srv.Close()

	// Delete the service from the registry
	err = srv.Delete()
	if err != nil {
		return err
	}

	return nil
}

// Status -
func (thisRef WindowsService) Status() (ServiceStatus, error) {
	status := ServiceStatus{}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("connecting to service manager: ", thisRef.command.Name),
	})

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error connecting to service manager: ", err),
		})
		return status, fmt.Errorf("could not connect to service manager: %v", err)
	}
	defer m.Disconnect()

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("opening system service"),
	})

	// Open the service so we can manage it
	srv, err := m.OpenService(thisRef.command.Name)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error opening service: ", err),
		})
		return status, fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	stat, err := srv.Query()
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("error getting service status: ", err),
		})
		return status, fmt.Errorf("could not get service status: %v", err)
	}

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprintf("service status: %#v", stat),
	})

	status.PID = int(stat.ProcessId)
	status.Running = stat.State == svc.Running

	return status, nil
}

// Exists -
func (thisRef WindowsService) Exists() bool {

	args := []string{"queryex", fmt.Sprintf("\"%s\"", thisRef.command.Name)}

	// https://www.computerhope.com/sc-command.htm
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("running command: sc ", strings.Join(args, " ")),
	})

	_, err := helpersExec.ExecWithArgs("sc", args...)
	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"method":  helpersReflect.GetThisFuncName(),
			"message": fmt.Sprint("exists service error: ", err),
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

func (thisRef WindowsService) control(command svc.Cmd, state svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(thisRef.command.Name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	status, err := srv.Control(command)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", command, err)
	}

	timeout := time.Now().Add(10 * time.Second)
	for status.State != state {
		// Exit if a timeout is reached
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", state)
		}

		time.Sleep(300 * time.Millisecond)

		// Make sure transition happens to the desired state
		status, err = srv.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}

	return nil
}

// Execute - implement the Windows `service.Handler` interface
func (thisRef WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"method":  helpersReflect.GetThisFuncName(),
		"message": fmt.Sprint("execute called"),
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
