package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type configFile struct {
	ShellConfigPath string            `json:"shellConfigPath,omitempty"`
	Aliases         map[string]string `json:"aliases,omitempty"`
	Apps            map[string]string `json:"apps,omitempty"`
	Paths           map[string]string `json:"paths,omitempty"`
}

func (cf *configFile) readConfig() error {
	// creating the dir
	dirPath, err := getConfigDirPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.sam dir")
		return err
	}

	//getting the file
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := getFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	//filling the data
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &cf)
	if err != nil {
		return err
	}

	if cf.Paths == nil {
		cf.Paths = make(map[string]string)
	}
	if cf.Apps == nil {
		cf.Apps = make(map[string]string)
	}
	if cf.Aliases == nil {
		cf.Aliases = make(map[string]string)
	}

	return nil
}

func (cf *configFile) writeConfig() error {
	//getting the file
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := getFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	//write into file
	content, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}

	file.Truncate(0)
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (cf *configFile) echo(items map[string]string) string {
	var buffer bytes.Buffer
	var counter = 0

	for key, value := range items {
		counter++
		buffer.WriteString(strconv.Itoa(counter) + ". " + key + " : " + value + "\n")
	}

	return buffer.String()
}

func (cf *configFile) ammend(apps map[string]string) error {
	for key, value := range apps {
		cf.Apps[key] = value
	}
	err := cf.writeConfig()
	return err
}

// get the apps inside ~/.sam
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
			var fileName = filepath.Base(path)
			if fileName == "config.json" {
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

// returns ~/.sam
func getConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sam"), nil
}

// returns ~/.sam/config.json
func getConfigFilePath() (string, error) {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, "config.json"), nil
}

func getReproduceFilePath() (string, error) {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, "reproduce.txt"), nil
}

// returns a file opened using rd|wr|create|0644 flags
func getFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// replaces/adds the content between the del in the file with content
func replaceFilePartition(del string, file *os.File, add bool, content ...string) error {
	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	var partitionedBuffer bytes.Buffer
	insideDel := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, del) && !insideDel {
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

// walks the baseDir (can be recursive) to returns all the files ending with 1 of the suffixes
func walkBaseDir(baseDir string, suffixes []string, recursive bool) (map[string]string, error) {
	aliases := make(map[string]string)
	if strings.HasPrefix(baseDir, "~/") {
		userDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		baseDir = filepath.Join(userDir, baseDir[2:])
	}

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
			aliases[name] = path
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

// copies a file into another dir
func copyFile(src string, dstDir string) (dstName string, err error) {
	newPath := strings.ToLower(filepath.Join(dstDir, strings.Replace(filepath.Base(src), " ", "-", -1)))
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

// moves a file into the ~/.sam dir
func storePath(src string) (dstName string, err error) {
	dst, err := getConfigDirPath()
	if err != nil {
		return "", err
	}

	return copyFile(src, dst)
}

func populateShellParser(cf configFile) ShellConfigParser {
	// ---------- gets shell parser ----------
	parser := getDynShellParser(cf)
	// ---------- gets shell parse ----------

	for k, v := range cf.Apps {
		parser.Add(Alias{k, v})
	}
	for k, v := range cf.Aliases {
		noSpace := strings.Replace(v, " ", "-", -1)
		app, exists := apps[noSpace]
		if exists {
			parser.Add(Alias{k, app})
		} else {
			parser.Add(Alias{k, v})
		}
	}
	for k, v := range cf.Paths {
		parser.ReproducePath(Alias{k, v})
	}

	return parser
}

func createReproduceFile(content []string) error {
	reproduceFilePath, err := getReproduceFilePath()
	if err != nil {
		return err
	}
	file, err := getFile(reproduceFilePath)
	defer file.Close()

	strContent := strings.Join(content, "\n")
	file.Truncate(0)
	_, err = file.Write([]byte(strContent))

	return nil
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
