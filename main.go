package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
		start:         true,
	}

	return model
}

type Model struct {
	inputField *textinput.Model
	filePicker *filepicker.Model

	pickingFile   bool
	typingCommand bool
	start         bool

	cmd string
}

type CommandOuput struct {
	output        string
	isOutput      bool
	isPickingFile bool
}

type StartFix struct {
}

func startFixMsg() tea.Msg {
	return StartFix{}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.filePicker.Init(), startFixMsg)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartFix:
		if m.start {
			go func() {
				time.Sleep(3 * time.Millisecond)
				m.typingCommand = true
				m.pickingFile = false
			}()
			return m, startFixMsg
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.pickingFile {
				m.typingCommand = true
				m.pickingFile = false
				return m, nil
			}
		case "enter":
			if m.pickingFile {
				break
			}
			if m.typingCommand {
				var cmd = m.handleCmd(m.inputField.Value())
				m.inputField.SetValue("")
				return m, cmd
			}
			return m, nil
		}
	case CommandOuput:
		if msg.isPickingFile {
			m.pickingFile = true
			m.typingCommand = false
			return m, nil
		}
		if msg.isOutput {
			return m, tea.Printf(msg.output)
		}
		return m, nil
	}
	if m.typingCommand {
		var inputField, cmd = m.inputField.Update(msg)
		m.inputField = &inputField
		return m, cmd
	}
	if m.pickingFile {
		var filePicker, cmd = m.filePicker.Update(msg)
		m.filePicker = &filePicker

		if didSelect, path := filePicker.DidSelectFile(msg); didSelect {

			if m.cmd == "fe" {
				var dir = filepath.Dir(path)
				os.Chdir(dir)
				cmd = tea.Batch(cmd, tea.Printf(dir))
			}
			if m.cmd == "ef" {
				command := exec.Command("code", path)
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
		m.start = false
		return fmt.Sprintf("CC %s\n", m.inputField.View())
	}

	if m.pickingFile {
		return fmt.Sprintf(m.filePicker.View())
	}

	return ""
}

func (m *Model) handleCmd(cmd string) tea.Cmd {
	return func() tea.Msg {
		if cmd == "fe" || cmd == "ef" {
			m.cmd = cmd
			return CommandOuput{"", false, true}
		}
		var output, isOutput = handleCommand(cmd)
		return CommandOuput{output, isOutput, false}
	}
}
