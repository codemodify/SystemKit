package tests

import (
	"github.com/codemodify/SystemKit/Service"
)

func createService() Service.SystemService {
	return Service.New(Service.ServiceCommand{
		Name:             "systemkit-test-service",
		DisplayLabel:     "SystemKit Test Service",
		Description:      "SystemKit Test Service",
		DocumentationURL: "",
		Executable:       "/usr/bin/vim",
		Args:             []string{""},
		WorkingDirectory: "/tmp",
		StdOutPath:       "null",
		RunAsUser:        "user",
	})
}
