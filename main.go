package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	var path = "C:\\Users\\meron\\Desktop\\cli-go\\Shortcuts\\Artix Game Launcher.lnk"
	var test = getShortcutCmd(path)

	fmt.Println(test)
	runLnk(test)

	//var initialModel = initialModel()

	//if _, err := tea.NewProgram(&initialModel, tea.WithAltScreen()).Run(); err != nil {
	//	panic(err)
	//}
}

func initialModel() Model {
	var model = Model{}
	var t = textinput.New()
	model.inputField = t

	return model
}

type Model struct {
	inputField textinput.Model

	Commands map[string]string
	Apps     map[string]string
	Aliases  map[string]string
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	return "Hello world!"
}
