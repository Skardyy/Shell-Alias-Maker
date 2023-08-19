package main

import (
	"fmt"

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
	var t = textinput.New()

	t.Focus()

	var model = Model{
		inputField: t,
		enter:      false,
	}

	return model
}

type Model struct {
	inputField textinput.Model
	enter      bool
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			go handleCommand(m.inputField.Value())
			m.enter = true
			m.inputField.Reset()
			return m, nil
		}
	}

	if m.enter {
		m.enter = false
		m.inputField.Cursor.Focus()
	}
	var cmd tea.Cmd
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return fmt.Sprintf("CC %s\n", m.inputField.View())
}
