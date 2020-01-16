# On OS X
- `launchd` is used
- service files is stored under `/Library/LaunchDaemons`
- helpers
    - `sudo launchctl stop SERVICE`
    - `sudo launchctl start SERVICE`
    - `sudo launchctl list | grep SERVICE`
    - `syslog -w`
