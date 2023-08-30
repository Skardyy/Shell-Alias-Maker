package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var aliases, shell = readConfig()
var shortcuts = getShortcuts()

// creates a cmd for the command, also gives a flag indicating to run async or not
func handleCommand(command string) (*exec.Cmd, bool) {
	if command == "?" {
		var cmd = echoCommands()
		return cmd, false
	}

	command = strings.ToLower(command)
	var alias, okAlias = aliases[command]
	if okAlias {
		var shortcut, okShortcut = shortcuts[alias.target]
		if okShortcut {
			var cmd = runLnk(shortcut)
			var flag = alias.t == "async"
			return cmd, flag
		}

		var cmd = exec.Command(shell, alias.target)
		var flag = alias.t == "async"

		return cmd, flag
	}

	var shortcut, okShortcut = shortcuts[command]
	if okShortcut {
		var cmd = runLnk(shortcut)
		return cmd, false
	}

	var cmd = exec.Command(shell, command)
	return cmd, false
}

// creates a cmd for a .lnk file
func runLnk(lnkCmd shortcut) *exec.Cmd {
	return exec.Command(shell, "&", fmt.Sprintf("'%s'", lnkCmd.target), lnkCmd.args)
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
	for key, value := range shortcuts {
		counter++
		buffer.WriteString(strconv.Itoa(counter) + ". " + key + " : " + value.target + " " + value.args + "\n")
	}

	return exec.Command(shell, "echo '"+buffer.String()+"'")
}
