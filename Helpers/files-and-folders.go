package Helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

// FileOrFolderExists -
func FileOrFolderExists(name string) bool {
	_, error := os.Stat(name)
	return !os.IsNotExist(error)
}

type FileChangedEventHandler func()

// WatchFile -
func WatchFile(filePath string, eventHandler FileChangedEventHandler) {
	if eventHandler == nil {
		log.Println(fmt.Sprintf("ERROR: %s",
			loggingC.Fields{
				"method": GetThisFuncName(),
				"error":  "File change handler is NULL",
			},
		))

		return
	}

	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s",
				loggingC.Fields{
					"method": GetThisFuncName(),
					"error":  err,
				},
			))

			initWG.Done()
		}
		defer watcher.Close()

		configFile := filepath.Clean(filePath)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filePath)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filePath)

					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {

						realConfigFile = currentConfigFile

						eventHandler()
					} else if filepath.Clean(event.Name) == configFile &&
						(event.Op&fsnotify.Remove&fsnotify.Remove) != 0 {

						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed

						log.Println(fmt.Sprintf("ERROR: %s",
							loggingC.Fields{
								"method": GetThisFuncName(),
								"error":  err,
							},
						))
					}
					eventsWG.Done()
					return
				}
			}
		}()
		watcher.Add(configDir)

		initWG.Done()
		eventsWG.Wait()
	}()

	initWG.Wait()
}
