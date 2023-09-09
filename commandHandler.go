package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var aliases, shell = readConfig()
var apps = getApps()

// creates a cmd for the command, also gives a flag indicating to run async or not
func handleCommand(command string) (*exec.Cmd, bool) {
	if command == "?" {
		var cmd = echoCommands()
		return cmd, false
	}

	command = strings.ToLower(command)
	var alias, okAlias = aliases[command]
	if okAlias {
		var app, okApp = apps[alias.target]
		if okApp {
			var cmd = runApp(app)
			var flag = alias.t == "async"
			return cmd, flag
		}

		var cmd = exec.Command(shell, alias.target)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		var flag = alias.t == "async"

		return cmd, flag
	}

	var app, okApp = apps[command]
	if okApp {
		var cmd = runApp(app)
		return cmd, false
	}

	var cmd = exec.Command(shell, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd, false
}

// creates a cmd for a .lnk file
func runApp(appTarget string) *exec.Cmd {
	var cmd = exec.Command(shell, ". ", fmt.Sprintf(`'%s'`, appTarget))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

// creates a cmd that echos all aliases and shortcuts
func echoCommands() *exec.Cmd {
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

	var cmd = exec.Command(shell, "echo '"+buffer.String()+"'")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
