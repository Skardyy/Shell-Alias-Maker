package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var aliases, shellConfigPath = readConfig()
var apps = getApps()

func echoAliases() string {
	var buffer bytes.Buffer
	var counter = 0

	for key, value := range aliases {
		counter++
		var str = strconv.Itoa(counter) + ". " + key + " : " + value + "\n"
		buffer.WriteString(str)
	}
	for key, value := range apps {
		counter++
		buffer.WriteString(strconv.Itoa(counter) + ". " + key + " : " + value + "\n")
	}

	return buffer.String()
}

func getApps() map[string]string {
	var apps map[string]string = make(map[string]string)

	dir := getConfigDirPath()

	var err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
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
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}

	return apps
}

// reads the config file to return the a: aliases and sc: shellConfig path
func readConfig() (a map[string]string, sc string) {

	file := getConfigFile()
	defer file.Close()

	var aliases map[string]string = make(map[string]string)
	var scanner = bufio.NewScanner(file)

	var shellConfigPath string
	if scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			shellConfigPath = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
		} else {
			panic("Can't read from config file with wrong formating")
		}
	}

	for scanner.Scan() {
		var byteSlice = scanner.Bytes()
		var line = bytes.NewBuffer(byteSlice).String()

		parts := strings.Split(line, ":")
		var length = len(parts)
		if length != 2 {
			// Skip lines that don't have the expected format
			continue
		}

		name := strings.ToLower(strings.TrimSpace(parts[0]))
		target := strings.ToLower(strings.TrimSpace(parts[1]))

		aliases[name] = target
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, ""
	}

	return aliases, shellConfigPath
}

// returns the path to cc config files
func getConfigDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".cc")
}
func getConfigFilePath() string {
	dirPath := getConfigDirPath()
	return filepath.Join(dirPath, "config.txt")
}

func createConfigDir() {
	dirPath := getConfigDirPath()

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.cc dir")
		panic(err)
	}
}

func getConfigFile() *os.File {
	createConfigDir()

	filePath := getConfigFilePath()
	return getFile(filePath)
}

func getFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func replaceFilePartition(del string, content []string, file *os.File) {
	scanner := bufio.NewScanner(file)
	var normalText []string
	insideDel := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, del) {
			insideDel = true
		} else if insideDel && strings.HasSuffix(line, del) {
			insideDel = false
		} else if !insideDel {
			normalText = append(normalText, line)
		}
	}

	normalText = append(normalText, content...)
	finalText := strings.Join(normalText, "\n")
	finalText = del + "\n" + finalText + "\n" + del
	clearFile(file)

	_, err := file.Write([]byte(finalText))
	if err != nil {
		panic(err)
	}
}

func clearFile(file *os.File) {
	err := file.Truncate(0)
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}
}

func walkBaseDir(baseDir string, suffixes []string, recursive bool) []Alias {

}
