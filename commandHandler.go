package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var aliases, shell = readConfig()
var apps = getApps()

// creates a cmd for the command
func handleCommand(command string) (cmd string, handled bool) {
	if command == "?" {
		var cmd = echoCommands()
		return cmd, true
	}

	command = strings.ToLower(command)
	var alias, okAlias = aliases[command]
	if okAlias {
		//app alias
		var app, okApp = apps[alias.target]
		if okApp {
			var flag = alias.t == "async"
			var cmd = runApp(app, flag)
			return cmd, true
		}

		//cmd alias
		var flag = alias.t == "async"
		var cmd = getRunner(flag) + alias.target
		return cmd, true
	}

	var app, okApp = apps[command]
	if okApp {
		var cmd = runApp(app, false)
		return cmd, true
	}

	//not handled
	return "", false
}

// creates a cmd for a .lnk file
func runApp(appTarget string, async bool) string {
	return getRunner(async) + fmt.Sprintf(`'%s'`, appTarget)
}

func getRunner(async bool) string {
	var runner string
	if async {
		runner = "& "
	} else {
		runner = ". "
	}
	return runner
}

// creates a cmd that echos all aliases and shortcuts
func echoCommands() string {
	var buffer bytes.Buffer
	var counter = 0

	for key, value := range aliases {
		counter++
		var str = strconv.Itoa(counter) + ". " + key + " : " + value.target
		if value.t != "" {
			str += " ! " + value.t + "\n"
		} else {
			str += "\n"
		}
		buffer.WriteString(str)
	}
	for key, value := range apps {
		counter++
		buffer.WriteString(strconv.Itoa(counter) + ". " + key + " : " + value + "\n")
	}

	return "echo '" + buffer.String() + "'" + "\n"
}
