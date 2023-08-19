package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var aliases = getAliases()
var shortcuts = getShortcuts()
var commands = getCommands()

func handleCommand(command string) (string, bool) {
	command = strings.ToLower(command)
	var alias, okAlias = aliases[command]
	if okAlias {
		var shortcut, okShortcut = shortcuts[alias]
		if okShortcut {
			runLnk(shortcut)
			return "Started " + alias, true
		}
	}

	var shortcut, okShortcut = shortcuts[command]
	if okShortcut {
		runLnk(shortcut)
		return "Started " + command, true
	}

	var cmd, okCmd = commands[command]
	if okCmd {
		cmd()
		return "", false
	}

	var res, _ = runCommand(command)
	return res, true
}

func runLnk(lnkCmd shortcut) {
	var cmd = exec.Command("powershell", "&", fmt.Sprintf("'%s'", lnkCmd.target), lnkCmd.args)
	cmd.Stdin = os.Stdin
	go cmd.Run()
}

func getCommands() map[string]func() {
	var cmds = make(map[string]func())

	cmds["fe"] = func() {
		//todo
	}
	cmds["ef"] = func() {
		//todo
	}

	return cmds
}

func runCommand(cmd string) (string, error) {
	command := exec.Command("powershell", cmd)
	command.Stdin = os.Stdin

	var stderr bytes.Buffer
	command.Stderr = &stderr

	var res, err = command.Output()
	if err != nil {
		return stderr.String(), err
	}
	return string(res), nil
}
