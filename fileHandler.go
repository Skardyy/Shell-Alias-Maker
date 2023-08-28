package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	lnk "github.com/parsiya/golnk"
)

type shortcut struct {
	target string
	args   string
}
type Alias struct {
	target string
	args   []string
	t      string
}

func getShortcut(path string) (shortcut, bool) {
	Lnk, err := lnk.File(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return shortcut{}, false
	}

	var cmd = Lnk.LinkInfo.LocalBasePath
	var args = Lnk.StringData.CommandLineArguments
	return shortcut{target: cmd, args: args}, true
}

// find a file with the given extension in the given root folder
func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func getShortcuts() map[string]shortcut {
	var shortcuts map[string]shortcut = make(map[string]shortcut)

	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var ShortcutFolder = fmt.Sprintf("%s\\Shortcuts", executablePath)

	for _, s := range find(ShortcutFolder, ".lnk") {
		var shortcut, ok = getShortcut(s)

		if !ok {
			continue
		}

		fileName := filepath.Base(s)
		extension := filepath.Ext(fileName)
		nameWithoutExtension := fileName[:len(fileName)-len(extension)]

		nameWithoutExtension = strings.ToLower(nameWithoutExtension)
		shortcuts[nameWithoutExtension] = shortcut
	}

	return shortcuts
}

func getAliases() map[string]Alias {
	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	var file, err = os.Open(fmt.Sprintf("%s\\Shortcuts\\aliases.txt", executablePath))
	if err != nil {
		log.Printf("Missing ~\\Shortcut\\aliases.txt")
		return nil
	}
	defer file.Close()

	var aliases map[string]Alias = make(map[string]Alias)

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {

		var line = scanner.Text()
		var parts = strings.Split(line, "!")
		var t string = ""
		if len(parts) == 2 {
			t = strings.ToLower(strings.TrimSpace(parts[1]))
		}

		parts = strings.Split(line, ":")
		var length = len(parts)
		if length < 2 {
			// Skip lines that don't have the expected format
			continue
		}
		name := strings.ToLower(strings.TrimSpace(parts[0]))
		target := strings.ToLower(strings.TrimSpace(parts[1]))

		var args []string
		if length > 2 {
			args = parts[2:length]
		}
		aliases[name] = Alias{target, args, t}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return aliases
}
