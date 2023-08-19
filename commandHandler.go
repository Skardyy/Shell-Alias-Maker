package main

import (
	"os"
	"os/exec"
	"strings"
)

var aliases = getAliases()
var shortcuts = getShortcuts()
var commands = getCommands()

// not tested
func handleCommand(command string) string {
	command = strings.ToLower(command)
	var alias, okAlias = aliases[command]
	if okAlias {
		var shortcut, okShortcut = shortcuts[alias]
		if okShortcut {
			runLnk(shortcut)
			return "Opened" + shortcut
		}
	}

	var shortcut, okShortcut = shortcuts[command]
	if okShortcut {
		runLnk(shortcut)
		return "Opened" + shortcut
	}

	var cmd, okCmd = commands[command]
	if okCmd {
		cmd()
		return ""
	}

	runCommand(command)
	return ""
}

// works
func runLnk(lnkCmd string) {
	var cmd = exec.Command(lnkCmd)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// works /not tested the cd part bcuz the app isnt done yet
func getCommands() map[string]func() {
	var cmds = make(map[string]func())

	cmds["fe"] = func() {
		runCommand("fzf | Split-Path | cd")
	}
	cmds["ef"] = func() {
		runCommand("fzf | % { code $_ }")
	}

	return cmds
}

// works
func runCommand(cmd string) error {
	command := exec.Command("powershell", "&", cmd)
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	var err = command.Run()
	if err != nil {
		return err
	}
	return nil
}
