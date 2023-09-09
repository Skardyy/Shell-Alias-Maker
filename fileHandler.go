package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Alias struct {
	target string
	t      string
}

func getApps() map[string]string {
	var apps map[string]string = make(map[string]string)

	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var appsFolder = fmt.Sprintf("%s\\Apps", executablePath)

	var err = filepath.Walk(appsFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var fileName = strings.ToLower(filepath.Base(path))
			if fileName == "readme.md" || fileName == "config.txt" {
				return nil
			}
			extension := filepath.Ext(fileName)
			nameWithoutExtension := fileName[:len(fileName)-len(extension)]
			apps[nameWithoutExtension] = path
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", appsFolder, err)
	}

	return apps
}

func readConfig() (map[string]Alias, string) {
	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var shell = "Powershell"

	var file, err = os.Open(fmt.Sprintf("%s\\Apps\\config.txt", executablePath))
	if err != nil {
		fmt.Println("Missing ~\\Apps\\config.txt")
		return nil, shell
	}
	defer file.Close()

	var aliases map[string]Alias = make(map[string]Alias)
	var scanner = bufio.NewScanner(file)

	if scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			shell = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
		}
	}

	for scanner.Scan() {
		var byteSlice = scanner.Bytes()
		var line = bytes.NewBuffer(byteSlice).String()
		var parts = strings.Split(line, "!")
		var t string = ""
		if len(parts) == 2 {
			t = strings.ToLower(strings.TrimSpace(parts[1]))
		}

		parts = strings.Split(parts[0], ":")
		var length = len(parts)
		if length != 2 {
			// Skip lines that don't have the expected format
			continue
		}

		name := strings.ToLower(strings.TrimSpace(parts[0]))
		var target string
		if t == "nolow" {
			target = strings.TrimSpace(parts[1])
		} else {
			target = strings.ToLower(strings.TrimSpace(parts[1]))
		}

		aliases[name] = Alias{target, t}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, shell
	}
	return aliases, shell
}
