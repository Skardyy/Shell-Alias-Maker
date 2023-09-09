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

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		var cmd, async = handleCommand(line)
		if async {
			go cmd.Run()
		} else {
			cmd.Run()
		}
	}
}
