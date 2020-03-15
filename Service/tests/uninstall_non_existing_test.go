package tests

import (
	"testing"

	helpersErrors "github.com/codemodify/SystemKit/Helpers"
	"github.com/codemodify/SystemKit/Service"
)

func Test_Uninstall(t *testing.T) {
	service := createService()

	err := service.Uninstall()
	if helpersErrors.Is(err, Service.ErrServiceDoesNotExist) {
		// INFO: this is a good thing
	} else if err != nil {
		t.Fatalf(err.Error())
	}
}
