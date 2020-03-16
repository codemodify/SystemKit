package tests

import (
	"fmt"

	helpersGuid "github.com/codemodify/SystemKit/Helpers"
	"github.com/codemodify/SystemKit/Service"
)

func createService() Service.SystemService {
	return Service.New(Service.Command{
		Name:             "systemkit-test-service",
		DisplayLabel:     "SystemKit Test Service",
		Description:      "SystemKit Test Service",
		DocumentationURL: "",
		Executable:       "/usr/bin/watch",
		Args:             []string{"ls"},
		WorkingDirectory: "/tmp",
		StdOutPath:       "null",
		RunAsUser:        "user",
	})
}

func createRandomService() Service.SystemService {
	randomData := helpersGuid.NewGUID()

	return Service.New(Service.Command{
		Name:             fmt.Sprintf("systemkit-test-service-%s", randomData),
		DisplayLabel:     fmt.Sprintf("SystemKit Test Service-%s", randomData),
		Description:      fmt.Sprintf("SystemKit Test Service-%s", randomData),
		DocumentationURL: "",
		Executable:       "/usr/bin/watch",
		Args:             []string{"ls"},
		WorkingDirectory: "/tmp",
		StdOutPath:       "null",
		RunAsUser:        "user",
	})
}
