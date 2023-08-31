package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	setLoggerOutput()
	var initialModel = initialModel()

	var dir, _ = os.UserHomeDir()
	os.Chdir(dir)

	if _, err := tea.NewProgram(&initialModel).Run(); err != nil {
		panic(err)
	}
}

func initialModel() Model {
	var t = createTextInput("?")

	var model = Model{
		inputField:    t,
		typingCommand: true,
	}

	return model
}

type Model struct {
	inputField *textinput.Model

	typingCommand bool

	cmd string
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.typingCommand {
				var cmd = m.handleCmd(m.inputField.Value())
				m.inputField.SetValue("")
				return m, cmd
			}
			return m, nil
		}
	}
	if m.typingCommand {
		var inputField, cmd = m.inputField.Update(msg)
		m.inputField = &inputField
		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	if m.typingCommand {
		return fmt.Sprintf("CC %s\n", m.inputField.View())
	}

	return ""
}

func (m *Model) handleCmd(cmd string) tea.Cmd {
	if cmd == "cc" {
		return nil
	}

	var command, async = handleCommand(cmd)
	if async {
		go command.Run()
		return nil
	}

	return tea.ExecProcess(command, nil)
}
