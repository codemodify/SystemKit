package tests

import (
	"testing"

	helpersErrors "github.com/codemodify/SystemKit/Helpers"
	"github.com/codemodify/SystemKit/Service"
)

func Test_stop(t *testing.T) {
	service := createService()

	err := service.Stop()
	if helpersErrors.Is(err, Service.ErrServiceDoesNotExist) {
		// this is a good thing
	} else if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_stop_non_existing(t *testing.T) {
	service := createRandomService()

	err := service.Stop()
	if helpersErrors.Is(err, Service.ErrServiceDoesNotExist) {
		// this is a good thing
	} else if err != nil {
		t.Fatalf(err.Error())
	}
}
