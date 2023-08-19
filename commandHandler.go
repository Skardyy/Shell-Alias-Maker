package main

import (
	"log"
	"os"
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

// works
func runLnk(lnkCmd string) {
	var cmd = exec.Command(lnkCmd)
	var err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
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
func runCommand(cmd string) {
	command := exec.Command("powershell", "&", cmd)
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
