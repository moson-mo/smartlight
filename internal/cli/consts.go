package cli

const (
	// info text which is printed when -h or equivalent parameters are used
	infoTXT = `smartlight cli-client

Usage:
smcli [COMMAND]
	
Commands:
start		starts the service

stop		stops the backlight service 
		(the service will still be running
		to accept further commands)

status		returns the current status

quit		quits the service
`
)
