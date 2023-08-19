package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	var initialModel = initialModel()

	if _, err := tea.NewProgram(&initialModel).Run(); err != nil {
		panic(err)
	}
}

func initialModel() Model {
	var t = createTextInput("Enter Command..")
	var f = createFilePicker()

	var model = Model{
		inputField:    t,
		filePicker:    f,
		pickingFile:   true,
		typingCommand: false,
	}

	return model
}

type Model struct {
	inputField textinput.Model
	filePicker filepicker.Model

	pickingFile   bool
	typingCommand bool

	cmd string
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.filePicker.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.pickingFile {
				break
			}
			var res, ok = m.handleCmd(m.inputField.Value())
			m.inputField.SetValue("")
			if ok {
				return m, tea.Printf(res)
			}
			return m, nil
		}
	}

	if m.typingCommand {
		var cmd tea.Cmd
		m.inputField, cmd = m.inputField.Update(msg)
		return m, cmd
	}
	if m.pickingFile {
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)

		if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {

			if m.cmd == "fe" {
				os.Chdir(filepath.Dir(path)) // would be better to let select folder, but seems to not work due to some bug, idk
			}
			if m.cmd == "ef" {
				command := exec.Command("code", path)
				go command.Run()
			}
			if m.cmd == "ef -a" {
				command := exec.Command("code", filepath.Dir(path))
				go command.Run()
			}

			m.typingCommand = true
			m.pickingFile = false
		}

		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	if m.typingCommand {
		return fmt.Sprintf("CC %s\n", m.inputField.View())
	}
	if m.pickingFile {
		return fmt.Sprintf(m.filePicker.View())
	}
	return ""
}

func (m *Model) handleCmd(cmd string) (string, bool) {
	if cmd == "fe" || cmd == "ef" || cmd == "ef -a" {
		m.cmd = cmd
		m.typingCommand = false
		m.pickingFile = true
		return "", false
	}
	return handleCommand(cmd)
}
