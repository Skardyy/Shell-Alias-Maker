package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

func main() {
	pr, pw := io.Pipe()
	cmd := exec.Command(shell)

	defer cmd.Wait()
	defer pw.Close()

	cmd.Stdin = pr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	//read stdin from the shell
	userScanner := bufio.NewScanner(os.Stdin)
	for userScanner.Scan() {
		input := userScanner.Text()

		cmd, handled := handleCommand(input)
		if !handled {
			pw.Write([]byte(input + "\n"))
		} else {
			pw.Write([]byte(cmd + "\n"))
		}
	}
}
