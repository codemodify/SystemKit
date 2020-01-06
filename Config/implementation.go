package Config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	helpersFile "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
)

var configInstance Config

// LoadConfig - loads a custom config type, takes an existing instance usually created with `&CustomAppConfig{}`
func LoadConfig(config Config) Config {
	configInstance = config.DefaultConfig()

	var configFileName = GetConfigFile()
	if helpersFile.FileOrFolderExists(configFileName) {
		file, err := ioutil.ReadFile(configFileName)
		if err != nil {
			logging.Instance().LogWarningWithFields(loggingC.Fields{
				"method": helpersReflect.GetThisFuncName(),
				"error":  fmt.Sprintf("unable to load config file %s, using default", configFileName),
			})
		} else {
			err := json.Unmarshal(file, configInstance)
			if err != nil {
				logging.Instance().LogWarningWithFields(loggingC.Fields{
					"method": helpersReflect.GetThisFuncName(),
					"error":  fmt.Sprintf("unable to load config file %s, using default", configFileName),
				})

				configInstance = config.DefaultConfig()
			}
		}
	}

	return configInstance
}

type ConfigChangedEventHandler func(config Config)

func OnConfigChanged(configChangedEventHandler ConfigChangedEventHandler) {
	helpersFile.WatchFile(GetConfigFile(), func() {
		configChangedEventHandler(LoadConfig(configInstance))
	})
}

// SaveConfig -
func SaveConfig(config Config) error {
	return ioutil.WriteFile(GetConfigFile(), []byte(config.String()), 0644)
}

// GetConfigDir -
func GetConfigDir() string {
	directoryName := "." + string(os.PathSeparator) + "config"
	if !helpersFile.FileOrFolderExists(directoryName) {
		err := os.MkdirAll(directoryName, 0755)
		if err != nil {
			logging.Instance().LogFatalWithFields(loggingC.Fields{
				"method": helpersReflect.GetThisFuncName(),
				"error":  fmt.Sprintf("unable to create directory %s", directoryName),
			})

			panic(err)
		}
	}
	return directoryName
}

// GetConfigFile -
func GetConfigFile() string {
	return GetConfigDir() + string(os.PathSeparator) + "config.json"
}
