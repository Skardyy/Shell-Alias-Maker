package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Alias struct {
	Name   string
	Target string
	App    bool
}

type configFile struct {
	Commands map[string]string `json:"commands"`
}

func initDirs() error {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.sam dir")
		return err
	}

	prePath := filepath.Join(dirPath, "pre")
	err = os.MkdirAll(prePath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.sam/pre dir")
		return err
	}

	dstPath := filepath.Join(dirPath, "dst")
	err = os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		fmt.Println("error creating ~/.sam dir")
		return err
	}

	return nil
}

func (cf *configFile) readConfig() error {
	// creating the dir
	err := initDirs()
	if err != nil {
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

	if len(content) != 0 {
		err = json.Unmarshal(content, &cf)
		if err != nil {
			return err
		}
	}

	if cf.Commands == nil {
		cf.Commands = make(map[string]string)
	}
	err = cf.writeConfig()
	if err != nil {
		return err
	}

	return nil
}

func (cf *configFile) formatPaths() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for key, value := range cf.Commands {
		if strings.Contains(strings.ToLower(value), strings.ToLower(homeDir)) {
			cf.Commands[key] = strings.Replace(value, homeDir, "~", 1)
		}
	}

	return nil
}

func (cf *configFile) deformatPaths() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for key, value := range cf.Commands {
		if strings.Contains(value, "~") {
			cf.Commands[key] = strings.Replace(value, "~", homeDir, 1)
		}
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

	//format the content
	err = cf.formatPaths()
	if err != nil {
		return err
	}
	//write into file
	content, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	//deformat the content just in case
	err = cf.deformatPaths()
	return err
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

// get the apps inside ~/.sam/pre
func getApps() (map[string]string, error) {
	var apps map[string]string = make(map[string]string)

	dir, err := getConfigPreDirPath()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var fileName = filepath.Base(path)
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

// returns ~/.sam/pre
func getConfigPreDirPath() (string, error) {
	dir, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "pre"), nil
}

// returns ~/.sam/dst
func getConfigDstDirPath() (string, error) {
	dir, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "dst"), nil
}

// returns ~/.sam/config.json
func getConfigFilePath() (string, error) {
	dirPath, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, "config.json"), nil
}

// returns a file opened using rd|wr|create|0644 flags
func getFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func promptTillSuccess(prompt string) string {
	var f string
	for true {
		color.Magenta(prompt)
		fmt.Scanln(&f)
		if f == "y" || f == "n" {
			return f
		}
	}
	return f
}

// writes the alias as a shell scripts so it can be ran afterwards
func parseAlias(alias Alias, ext string, cf *configFile) error {
	path, err := getConfigDstDirPath()
	if err != nil {
		return err
	}

	file, err := getFile(filepath.Join(path, alias.Name+"."+ext))
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if len(content) == 0 {
		_, err = file.WriteString(alias.Target)
		color.Green("writing " + alias.Name + "." + ext)
		if err != nil {
			fmt.Println("had an issue parsing ", alias.Name, alias.Target)
			fmt.Println(err)
		}
	} else {
		strContent := string(content)
		if strContent != alias.Target {
			fmt.Println(alias.Name+"."+ext, "had an conflict, should i take the one at config.json?")
			fmt.Println("Original Content:")
			color.Cyan(strContent)
			fmt.Println("config.json Content:")
			color.Cyan(alias.Target)
			f := promptTillSuccess("replace its content with the one from the config.json? y/n")
			if f == "n" {
				cf.Commands[alias.Name] = strContent
			} else {
				file.Truncate(0)
				file.Seek(0, 0)
				_, err = file.WriteString(alias.Target)
				color.Green("writing " + alias.Name + "." + ext)
				if err != nil {
					fmt.Println("had an issue parsing ", alias.Name, alias.Target)
					fmt.Println(err)
				}
			}
		}
	}

	return err
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

// moves a file into the ~/.sam/pre dir
func storePath(src string) (dstName string, err error) {
	dst, err := getConfigPreDirPath()
	if err != nil {
		return "", err
	}

	color.Yellow("Copying " + src)
	return copyFile(src, dst)
}
