package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var aliases = getAliases()
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

		var args = []string{alias.target}
		args = append(args, alias.args...)
		var cmd = exec.Command("Powershell", args...)

		var flag = alias.t == "async"
		return cmd, flag
	}

	var shortcut, okShortcut = shortcuts[command]
	if okShortcut {
		var cmd = runLnk(shortcut)
		return cmd, false
	}

	var cmd = exec.Command("powershell", command)
	return cmd, false
}

// creates a cmd for a .lnk file
func runLnk(lnkCmd shortcut) *exec.Cmd {
	return exec.Command("powershell", "&", fmt.Sprintf("'%s'", lnkCmd.target), lnkCmd.args)
}

// creates a cmd that echos all aliases and shortcuts
func echoCommands() *exec.Cmd {
	var buffer bytes.Buffer
	var counter = 0

	for key, value := range aliases {
		counter++
		var str = strconv.Itoa(counter) + ". " + key + " : " + value.target

		var args = argsToString(value.args)
		str += args

		if value.t == "async" {
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

	return exec.Command("powershell", "echo '"+buffer.String()+"'")
}

func argsToString(args []string) string {
	var buffer bytes.Buffer
	for _, arg := range args {
		buffer.WriteString(arg + " ")
	}
	return buffer.String()
}
