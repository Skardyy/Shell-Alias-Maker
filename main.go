package main

import (
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

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
