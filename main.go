package main

import (
	"os"

	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var dir, _ = os.UserHomeDir()
	os.Chdir(dir)

	commands := []string{}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		if len(line) > 0 {
			commands = append(commands, line)
		}

		var cmd, async = handleCommand(line)
		if async {
			go cmd.Run()
		} else {
			cmd.Run()
		}
	}
}
