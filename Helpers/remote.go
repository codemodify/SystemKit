package Helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sfreiberg/simplessh"
)

// NewSSHClient -
func NewSSHClient(host string, user string, pemFilePath string) (*simplessh.Client, error) {
	privateKeyAsString, err := ReadFileAsString(pemFilePath)
	if err != nil {
		return nil, err
	}

	return simplessh.ConnectWithKey(host, user, privateKeyAsString)
}

// RemoteUploadFiles -
func RemoteUploadFiles(host string, user string, pemFilePath, remoteDir string, localFiles []string) error {
	client, err := NewSSHClient(host, user, pemFilePath)
	if err != nil {
		return err
	}

	exists, err := RemoteExists(client, remoteDir)
	if err != nil {
		return err
	}

	if !exists {
		_, err := client.Exec(fmt.Sprintf("sudo mkdir -m 777 -p %s\n", remoteDir))
		if err != nil {
			return err
		}
	}

	for _, file := range localFiles {
		fName, err := os.Stat(file)
		if err != nil {
			return err
		}
		remoteFile := fmt.Sprintf("%s/%s", remoteDir, fName.Name())

		err = client.Upload(file, remoteFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoteDownload -
func RemoteDownload(host string, user string, pemFilePath string, localDir string, remoteFiles []string) error {
	client, err := NewSSHClient(host, user, pemFilePath)
	if err != nil {
		return err
	}

	if !FileOrFolderExists(localDir) {
		MakeDirTree(localDir)
	}

	for _, file := range remoteFiles {
		fName := filepath.Base(file)
		if err != nil {
			return err
		}
		localFile := fmt.Sprintf("%s/%s", localDir, fName)
		err = client.Download(file, localFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoteExists -
func RemoteExists(client *simplessh.Client, remoteDir string) (bool, error) {
	check := fmt.Sprintf("test -d %s && echo 'true' || echo 'false'", remoteDir)
	result, err := client.Exec(check)
	if err != nil {
		return false, err
	}
	trim := strings.TrimSuffix(string(result), "\n")
	return strconv.ParseBool(trim)
}
