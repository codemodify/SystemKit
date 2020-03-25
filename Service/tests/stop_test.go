package tests

import (
	"testing"

	"github.com/codemodify/SystemKit/Service"
	helpersErrors "github.com/codemodify/systemkit-helpers"
)

func Test_stop(t *testing.T) {
	service := CreateRemoteitService()

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
