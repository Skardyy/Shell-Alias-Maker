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

	var res, _ = runCommand(command)
	return res, true
}

func runLnk(lnkCmd shortcut) {
	var cmd = exec.Command("powershell", "&", fmt.Sprintf("'%s'", lnkCmd.target), lnkCmd.args)
	cmd.Stdin = os.Stdin
	go cmd.Run()
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
