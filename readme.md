# ![](https://fonts.gstatic.com/s/i/materialicons/label_important/v4/24px.svg) SystemKit
[![GoDoc](https://godoc.org/github.com/codemodify/SystemKit?status.svg)](https://godoc.org/github.com/codemodify/SystemKit)
[![0-License](https://img.shields.io/badge/license-0--license-brightgreen)](https://github.com/codemodify/TheFreeLicense)
[![Go Report Card](https://goreportcard.com/badge/github.com/codemodify/SystemKit)](https://goreportcard.com/report/github.com/codemodify/SystemKit)
[![Test Status](https://github.com/danawoodman/systemservice/workflows/Test/badge.svg)](https://github.com/danawoodman/systemservice/actions)
![code size](https://img.shields.io/github/languages/code-size/codemodify/SystemKit?style=flat-square)

Is a set of carefully cherry picked approaches to have stable framework to write complex applications fast

See samples or readme for each folder

# ![](https://fonts.gstatic.com/s/i/materialicons/label_important/v4/24px.svg) Runs On

<nobr>
<img src="https://img.icons8.com/ios-filled/50/000000/linux.png" width="30" />
<nobr /> <img src="https://img.icons8.com/ios-filled/50/000000/raspberry-pi.png" width="30" />
<nobr /> <img src="https://img.icons8.com/ios-filled/50/000000/mac-os.png" width="30" />
<nobr /> <img src="https://img.icons8.com/ios-filled/50/000000/windows-logo.png" width="30" />
</nobr>

# ![](https://fonts.gstatic.com/s/i/materialicons/label_important/v4/24px.svg) SDLC
- Debug using VSCode as `root`
	- `sudo dlv debug --headless --listen=:2345 --log --api-version=2 -- ANY-ARGS-HERE`
	- Attach to debugger using
		```json
		{
			"name": "Attach to SUDO debbugger",
			"type": "go",
			"request": "launch",
			"mode": "remote",
			"program": "${workspaceFolder}/THE-MAIN.go",
			"port": 2345,
			"host": "127.0.0.1",
			"remotePath": "${workspaceFolder}/THE-MAIN.go"
		}
		```

- Test on other platforms
	- `export VAGRANT_VAGRANTFILE=.helper-files/Vagrantfile`
	- `export VAGRANT_DOTFILE_PATH=.helper-files/.vagrant`
	- `vagrant up windows` (see ./.helper-files/Vagrantfile for more platforms)
	- `vagrant ssh windows` <- do the testing
	- `vagrant halt windows` <- power off the VM
	- `vagrant destroy windows -f` <- destroy the VM

- Install Go
	- windows
		- `cd \Users\vagrant\Downloads`
		-
		- Busybox for `ls`, `vi`, `more`, `grep`, `find` and such
			- `curl -O https://frippery.org/files/busybox/busybox64.exe`
			- `mkdir \busybox`
			- `busybox64.exe --install \busybox`
			- `setx path "%path%;\busybox"`
			- OR
			- `\Users\vagrant\Downloads\busybox64.exe sh` to have the `Ctrl + L`
		-
		- `curl -O https://dl.google.com/go/go1.13.7.windows-amd64.msi`
		- `msiexec /i go1.13.7.windows-amd64.msi /quiet /qn /norestart /log install.log`
		- `shutdown -s -t 0`
		-
		- `vagrant up windows && vagrant ssh windows`
		- `cd /vagrant/Service/tests`
		- `sc query | more` <- pick a service, for example `Spooler`
		- `vi sample-01_test.go` <- set the service name from above to test with
		- `go test -run Sample01_Status`
		- `sc stop Spooler`
		- `go test -run Sample01_Status`
		- `sc start Spooler`
		- `go test -run Sample01_Status`