package tests

import (
	"fmt"
	"testing"

	logging "github.com/codemodify/SystemKit/Logging"
)

func Test_LoadConfig(t *testing.T) {
	logging.Init(logging.NewConsoleLogger())

	config := getConfig()

	t.Log(config)
	fmt.Println(config)
}

func getConfig() *AppConfig {
	appDefaultConfig := &AppConfig{}
	return appDefaultConfig.CreateOrReturnInstance().(*AppConfig)
}
