package main

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func createTextInput(placeholder string) textinput.Model {
	var t = textinput.New()
	t.Placeholder = placeholder
	t.Focus()

	return t
}
