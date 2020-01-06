package Config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	logging "github.com/codemodify/SystemKit/Logging"
	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"

	helpersFile "github.com/codemodify/SystemKit/Helpers"
	helpersReflect "github.com/codemodify/SystemKit/Helpers"
)

var configInstance Config
var configOnce sync.Once

// LoadConfig -
func LoadConfig(config Config) Config {
	configOnce.Do(func() {
		configInstance = config.DefaultConfig()

		var configFileName = GetConfigDir() + string(os.PathSeparator) + "config.json"
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
	})

	return configInstance
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
