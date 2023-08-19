package main

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
)

func createTextInput(placeholder string) textinput.Model {
	var t = textinput.New()
	t.Placeholder = placeholder
	t.Focus()

	return t
}

func createFilePicker() filepicker.Model {
	var f = filepicker.New()
	f.CurrentDirectory, _ = os.Getwd()
	f.DirAllowed = false // disabled to some selecting bug idk
	f.FileAllowed = true
	return f
}
