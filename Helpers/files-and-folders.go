package Helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

// FileOrFolderExists -
func FileOrFolderExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// IsFolder -
func IsFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode().IsDir()
}

// IsFile -
func IsFile(path string) bool {
	return !IsFolder(path)
}

// MakeDirTree -
func MakeDirTree(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	return nil
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

// ReadFileAsBytes -
func ReadFileAsBytes(path string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

// ReadFileAsString -
func ReadFileAsString(path string) (string, error) {
	bytes, err := ReadFileAsBytes(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ReadFileAsLines -
func ReadFileAsLines(path string) ([]string, error) {
	fileAsString, err := ReadFileAsString(path)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(fileAsString, "\n"), nil
}

// ReadFileAsCSV -
func ReadFileAsCSV(path string) ([][]string, error) {
	lines, err := ReadFileAsLines(path)
	if err != nil {
		return [][]string{}, err
	}

	result := [][]string{}
	for _, line := range lines {
		result = append(result, strings.Split(line, ";"))
	}

	return result, nil
}

// ListFolderContent -
func ListFolderContent(folder string, ext string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			if !IsNullOrEmpty(ext) {
				if filepath.Ext(info.Name()) == ext {
					fileList = append(fileList, path)
				}
			} else {
				fileList = append(fileList, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
