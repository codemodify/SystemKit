package main

import (
	"testing"

	service "github.com/codemodify/SystemKit/Service"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

var cmd service.ServiceCommand

func init() {
	cmd = service.ServiceCommand{
		Name:         "Spooler",
		DisplayLabel: "My Service",
		Executable:   "vim",
		Args:         []string{},
		Description:  "My systemservice test!",
	}
}

func Test_Sample01_Run(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"run": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Run()
}

func Test_Sample01_Install(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"install": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Install(false)
}

func Test_Sample01_Start(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"start": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Start()
}

func Test_Sample01_Restart(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"restart": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Restart()
}

func Test_Sample01_Stop(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"stop": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Stop()
}

func Test_Sample01_Uninstall(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"uninstall": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Uninstall()
}

func Test_Sample01_Status(t *testing.T) {
	syetemService := service.New(cmd)
	status, err := syetemService.Status()

	if err != nil {
		logging.Instance().LogErrorWithFields(loggingC.Fields{
			"err": err.Error(),
		})
	} else {
		logging.Instance().LogInfoWithFields(loggingC.Fields{
			"Running": status.IsRunning,
			"PID":     status.PID,
		})
	}
}

func Test_Sample01_Exists(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	logging.Instance().LogInfoWithFields(loggingC.Fields{
		"exists": cmd.String(),
	})

	syetemService := service.New(cmd)
	syetemService.Exists()
}
