package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var aliases, shellConfigPath, _ = readConfig()
var apps, _ = getApps()

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

func getApps() (map[string]string, error) {
	var apps map[string]string = make(map[string]string)

	dir, err := getConfigDirPath()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var fileName = strings.ToLower(filepath.Base(path))
			if fileName == "config.txt" {
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

	return apps, nil
}

// reads the config file to return the a: aliases and sc: shellConfig path
func readConfig() (a map[string]string, sc string, err error) {

	file, err := getConfigFile()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	var aliases map[string]string = make(map[string]string)
	var scanner = bufio.NewScanner(file)

	var shellConfigPath string
	correctFormat := false
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			shellConfigPath = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			correctFormat = true
			break
		}
	}
	if !correctFormat {
		panic("Can't read from config file with wrong formating")
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
		return nil, "", err
	}

	return aliases, shellConfigPath, nil
}

// returns the path to cc config files
func getConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cc"), nil
}
func getConfigFilePath() (string, error) {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, "config.txt"), nil
}

func createConfigDir() error {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.cc dir")
		return err
	}

	return nil
}

func getConfigFile() (*os.File, error) {
	createConfigDir()

	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	return getFile(filePath)
}

func getFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func replaceFilePartition(del string, file *os.File, add bool, content ...string) error {
	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	var partitionedBuffer bytes.Buffer
	insideDel := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, del) {
			insideDel = true
		} else if insideDel && strings.HasPrefix(line, del) {
			insideDel = false
		} else if !insideDel {
			buffer.WriteString(line + "\n")
		} else {
			partitionedBuffer.WriteString(line + "\n")
		}
	}

	if !add {
		partitionedBuffer.Truncate(0)
	}
	partitionedBuffer.WriteString(strings.Join(content, "\n"))
	finalContent := checkDup(partitionedBuffer)
	buffer.WriteString(del + "\n" + finalContent + "\n" + del)
	file.Close()
	os.WriteFile(file.Name(), buffer.Bytes(), os.ModePerm)

	return nil
}

func walkBaseDir(baseDir string, suffixes []string, recursive bool) ([]Alias, error) {
	aliases := make([]Alias, 0)

	err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !recursive {
			if path != baseDir {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && containsExt(filepath.Ext(path), suffixes) {
			name := filepath.Base(path)
			ext := filepath.Ext(name)
			name = name[:len(name)-len(ext)]
			aliases = append(aliases, Alias{name, path})
		}
		return nil
	})

	return aliases, err
}

func containsExt(ext string, exts []string) bool {
	for _, v := range exts {
		if v == ext {
			return true
		}
	}
	return false
}

// better use replace file partition as its more delicate
func changeConfigFileContent(newShellConfigPath string, content string, add bool, file *os.File) error {
	var buffer bytes.Buffer

	scanner := bufio.NewScanner(file)

	//shell config file changing
	if scanner.Scan() {
		line := scanner.Text()
		if newShellConfigPath != "" {
			buffer.WriteString("[" + newShellConfigPath + "]" + "\n")
		} else {
			buffer.WriteString(line + "\n")
		}
	} else {
		buffer.WriteString("[" + newShellConfigPath + "]" + "\n")
	}

	//replacing the content of the file
	if content == "" || add {
		for scanner.Scan() {
			line := scanner.Text()
			buffer.WriteString(line + "\n")
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	if add || content != "" {
		buffer.WriteString(content)
	}

	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	checkedString := checkDup(buffer)
	err = os.WriteFile(configFilePath, []byte(checkedString), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func amendApps() error {
	file, err := getConfigFile()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	for k, v := range apps {
		buffer.WriteString(k + " : " + v + "\n")
	}

	return replaceFilePartition("#Apps", file, false, buffer.String())
}

func checkDup(buffer bytes.Buffer) string {
	slice := make([]string, 0)
	tokens := strings.Split(buffer.String(), "\n")
	for _, v := range tokens {
		if slices.Contains(slice, v) {
			if !strings.HasPrefix(v, "#") {
				continue
			}
		}
		if v != "" {
			slice = append(slice, v)
		}
	}
	return strings.Join(slice, "\n")
}

func getPopulatedShellParser() ShellConfigParser {
	parser := ShellConfigParser{}
	parser.With(shellConfigPath, &PwshConfigParsser{})

	for k, v := range apps {
		parser.Add(Alias{k, v})
	}
	for k, v := range aliases {
		noSpace := strings.Replace(v, " ", "-", -1)
		app, exists := apps[noSpace]
		if exists {
			parser.Add(Alias{k, app})
		} else {
			parser.Add(Alias{k, v})
		}
	}

	return parser
}

func storePath(src string) (dstName string, err error) {
	dst, err := getConfigDirPath()
	if err != nil {
		return "", err
	}

	return copyFile(src, dst)
}

func copyFile(src string, dstDir string) (dstName string, err error) {
	newPath := filepath.Join(dstDir, strings.Replace(filepath.Base(src), " ", "-", -1))
	srcFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", err
	}

	return newPath, dstFile.Sync()
}
