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

// works
func getShortcutCmd(path string) string {
	Lnk, err := lnk.File(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}

	var cmd = Lnk.LinkInfo.LocalBasePath
	var args = Lnk.StringData.CommandLineArguments
	cmd += " " + args
	return cmd
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

// works
func getShortcuts() map[string]string {
	var shortcuts map[string]string = make(map[string]string)

	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var ShortcutFolder = fmt.Sprintf("%s\\Shortcuts", executablePath)

	for _, s := range find(ShortcutFolder, ".lnk") {
		var target = getShortcutCmd(s)

		fileName := filepath.Base(s)
		extension := filepath.Ext(fileName)
		nameWithoutExtension := fileName[:len(fileName)-len(extension)]

		shortcuts[nameWithoutExtension] = target
	}

	return shortcuts
}

// works
func getAliases() map[string]string {
	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	file, err := os.Open(fmt.Sprintf("%s\\Shortcuts\\aliases.txt", executablePath))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	var shortcuts map[string]string = make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			// Skip lines that don't have the expected format
			continue
		}
		name := strings.TrimSpace(parts[0])
		target := strings.TrimSpace(parts[1])
		shortcuts[name] = target
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return shortcuts
}
