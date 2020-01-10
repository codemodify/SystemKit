// +build windows

package Service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

/*
Run is the process which gets fired when the service starts up
when the service is installed and started.
*/
func (s *SystemService) Run() error {
	logger.Log("running service")

	name := s.Command.Name
	debugOn := s.Command.Debug

	var err error
	if debugOn {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			logger.Log("error opening logs: ", err)
			return err
		}
	}
	defer elog.Close()

	logger.Log("starting service: ", name)
	elog.Info(1, fmt.Sprintf("starting %s service", name))

	run := svc.Run
	if debugOn {
		run = debug.Run
	}

	err = run(name, &windowsService{})
	if err != nil {
		logger.Log("error running service: ", err)
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return err
	}

	logger.Log("service stopped: ", name)
	elog.Info(1, fmt.Sprintf("%s service stopped", name))

	// if err := svc.Run(s.Command.Name, &windowsService{}); err != nil {
	// 	return err
	// }

	return nil
}

/*
Install the system service. If start is passed, also starts
the service.
*/
func (s *SystemService) Install(start bool) error {
	name := s.Command.Name
	exePath := s.Command.Program
	args := s.Command.Args
	desc := s.Command.Description

	logger.Log("installing system service: ", name)

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logger.Log("error connecting to service manager: ", err)
		return err
	}
	defer m.Disconnect()

	// Open the service so we can manage it
	srv, err := m.OpenService(name)
	if err == nil {
		logger.Log("error opening the service: ", name)
		srv.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	logger.Logf("creating service \"%s\" at path \"%s\" with args \"%s\"", name, exePath, args)

	// Create the system service
	conf := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: name,
		Description: desc,
	}
	srv, err = m.CreateService(name, exePath, conf, args...)
	if err != nil {
		logger.Log("error creating service: ", err)
		return err
	}
	defer srv.Close()

	// Remove event log if it is there
	_ = eventlog.Remove(name)

	logger.Log("setting up event logs: ", name)

	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		logger.Log("error creating service logs: ", err)
		srv.Delete()
		return fmt.Errorf("setting up event log failed: %s", err)
	}

	logger.Log("starting service: ", name)
	if start {
		if err := s.Start(); err != nil {
			logger.Log("error starting service: ", err)
			return err
		}
	}

	return nil

	// logger.Log("install service")

	// name := s.Command.Name
	// prog := s.Command.String()
	// args := []string{
	// 	"create",
	// 	fmt.Sprintf("\"%s\"", name),
	// 	"binPath=",
	// 	fmt.Sprintf("\"%s\"", prog),
	// 	// "start=",
	// 	// "boot",
	// }

	// out, err := runScCommand(args...)

	// if err != nil {
	// 	if strings.Contains(err.Error(), "exit status 1073") {
	// 		logger.Log("service already exists")
	// 	} else {
	// 		logger.Log("sc create output:\n", out)
	// 		return err
	// 	}
	// }

	// // if strings.Contains(out, "SUCCESS") {
	// // 	return nil
	// // }
}

/*
Start the system service if it is installed
*/
func (s *SystemService) Start() error {
	name := s.Command.Name

	logger.Log("starting system service: ", name)

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logger.Log("error connecting to service manager: ", err)
		return fmt.Errorf("could not connect to service manager: %v", err)
	}
	defer m.Disconnect()

	logger.Log("opening system service")

	// Open the service so we can manage it
	srv, err := m.OpenService(name)
	if err != nil {
		logger.Log("error opening service: ", err)
		return fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	logger.Log("attempting to start system service")

	err = srv.Start(s.Command.Args...)
	if err != nil {
		logger.Log("error starting service: ", err)
		return fmt.Errorf("could not start service: %v", err)
	}

	logger.Log("running service")

	return nil
	// _, err := runScCommand("start", fmt.Sprintf("\"%s\"", s.Command.Name))

	// if err != nil {
	// 	logger.Log("start service error: ", err)
	// 	return err
	// }

	// return nil
}

/*
Restart attempts to stop the service if running then starts it again
*/
func (s *SystemService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}

	if err := s.Start(); err != nil {
		return err
	}

	return nil
}

/*
Stop stops the system service by unloading the unit file
*/
func (s *SystemService) Stop() error {
	err := s.control(svc.Stop, svc.Stopped)
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

		logger.Log("waiting for service to stop")

		// // Wait a few seconds before retrying.
		time.Sleep(wait)

		// // Attempt to start the service again.
		stat, err := s.Status()
		if err != nil {
			return err
		}

		// // Check the status to see if it is running yet.
		// stat, err := system.Service.Status()
		// if err != nil {
		// 	exit(err, stop)
		// }

		// // If it is now running, exit the retry loop.
		if !stat.Running {
			break
		}

		if attempt == maxAttempts {
			return errors.New("could not stop system service after multiple attempts")
		}
	}

	return nil
	// _, err := runScCommand("stop", fmt.Sprintf("\"%s\"", s.Command.Name))

	// if err != nil {
	// 	logger.Log("stop service error: ", err)

	// 	if strings.Contains(err.Error(), "exit status 1062") {
	// 		logger.Log("service already stopped")
	// 	} else {
	// 		return err
	// 	}
	// }

	// return nil
}

/*
Uninstall the system service by first stopping it then removing
the unit file.
*/
func (s *SystemService) Uninstall() error {
	name := s.Command.Name

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// Open the service so we can manage it
	srv, err := m.OpenService(name)
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

	// Remove the event log
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("removing event log failed: %s", err)
	}

	return nil
	// name := s.Command.Name

	// err := s.Stop()

	// if err != nil {
	// 	return err
	// }

	// _, err = runScCommand("delete", fmt.Sprintf("\"%s\"", name))

	// if err != nil {
	// 	logger.Log("delete service error: ", err)
	// 	return err
	// }

	// return nil
}

/*
Status returns whether or not the system service is running
*/
func (s *SystemService) Status() (status *ServiceStatus, err error) {
	name := s.Command.Name
	status = &ServiceStatus{}

	logger.Log("connecting to service manager: ", name)

	// Connect to Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		logger.Log("error connecting to service manager: ", err)
		return status, fmt.Errorf("could not connect to service manager: %v", err)
	}
	defer m.Disconnect()

	logger.Log("opening system service")

	// Open the service so we can manage it
	srv, err := m.OpenService(name)
	if err != nil {
		logger.Log("error opening service: ", err)
		return status, fmt.Errorf("could not access service: %v", err)
	}
	defer srv.Close()

	stat, err := srv.Query()
	if err != nil {
		logger.Log("error getting service status: ", err)
		return status, fmt.Errorf("could not get service status: %v", err)
	}

	logger.Logf("service status: %#v", stat)

	status.PID = int(stat.ProcessId)
	status.Running = stat.State == svc.Running
	return status, nil
}

/*
Return whether or not the unit file eixts
*/
func (s *SystemService) Exists() bool {
	_, err := runScCommand("queryex", fmt.Sprintf("\"%s\"", s.Command.Name))

	if err != nil {
		logger.Log("exists service error: ", err)
		// Service does not exist
		// if strings.Contains(err.Error(), "FAILED 1060") {
		// return false
		// }
		return false
	}

	return true
}

func (s *SystemService) control(command svc.Cmd, state svc.State) error {
	name := s.Command.Name

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	srv, err := m.OpenService(name)
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

// var beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")

// func beep() {
// 	beepFunc.Call(0xffffffff)
// }

// logger.Log("uninstall service: ", name)

// serv, err := connectService(name)

// if err != nil {
// 	logger.Log("uninstall error: ", err)
// 	return err
// }

// defer serv.Close()

// logger.Log("deleting service")

// err = serv.Delete()

// if err != nil {
// 	logger.Log("delete service error: ", err)
// 	return err
// }

// logger.Log("service deleted")
// logger.Log("removing event log")

// err = eventlog.Remove(name)

// if err != nil {
// 	logger.Log("remove event log error: ", err)
// 	return err
// }

// logger.Log("event log removed")

// manager, err := openManager()

// if err != nil {
// 	return err
// }

// defer manager.Disconnect()

// cmd := s.Command

// serv, err := manager.CreateService(
// 	cmd.Name,
// 	cmd.Program,
// 	mgr.Config{DisplayName: cmd.Name},
// 	cmd.Args...,
// )

// if err != nil {
// 	return err
// }

// defer serv.Close()

// err = eventlog.InstallAsEventCreate(cmd.Name, eventlog.Error|eventlog.Warning|eventlog.Info)

// if err != nil {
// 	serv.Delete()
// 	return fmt.Errorf("SetupEventLogSource() failed: %s", err)
// }

// logger.Logf("manager: %+v", manager)

// serv, err := connectService(name)

// if err != nil {
// 	logger.Log("status error: ", err)
// 	return ServiceStatus{}, err
// }

// defer serv.Close()

// logger.Logf("service: %+v", serv)

// pid := getPID(name)
