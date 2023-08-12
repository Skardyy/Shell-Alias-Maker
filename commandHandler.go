package main

import (
	"log"
	"os/exec"
)

var aliases = getAliases()
var shortcuts = getShortcuts()
var commands = getCommands()

// not tested
func handleCommand(command string) {
	var alias, okAlias = aliases[command]
	if okAlias {
		var shortcut, okShortcut = shortcuts[alias]
		if okShortcut {
			runLnk(shortcut)
			return
		}
	}

	var shortcut, okShortcut = shortcuts[command]
	if okShortcut {
		runLnk(shortcut)
		return
	}

	var cmd, okCmd = commands[command]
	if okCmd {
		cmd()
		return
	}
}

// not tested
func runLnk(lnkPath string) {
	var cmd = exec.Command("powershell", "&", lnkPath)
	var err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// not tested
func getCommands() map[string]func() {
	var cmds = make(map[string]func())

	cmds["fe"] = func() {
		runCommand("Powershell fzf | Split-Path | cd")
	}
	cmds["ef"] = func() {
		runCommand("Powershell fzf | $ { code $_ }")
	}

	return cmds
}

// not tested
func runCommand(cmd string) {
	var command = exec.Command(cmd)
	var err = command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
