package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ShellLinkHeader struct {
	HeaderSize [4]byte  //HeaderSize
	ClassID    [16]byte //LinkCLSID
	LinkFlags  uint32   //LinkFlags      [4]byte
	FileAttr   uint32   //FileAttributes [4]byte
	Creation   [8]byte  //CreationTime
	Access     [8]byte  //AccessTime
	Write      [8]byte  //WriteTime
	FileSz     [4]byte  //FileSize
	IconIndex  [4]byte  //IconIndex
	ShowCmd    [4]byte  //ShowCommand

	//[2]byte HotKey values for shortcut shortcuts
	HotKeyLow  byte //HotKeyLow
	HotKeyHigh byte //HotKeyHigh

	Reserved1 [2]byte //Reserved1
	Reserved2 [4]byte //Reserved2
	Reserved3 [4]byte //Reserved3
}
type Shortcut struct {
	Name   string
	Target string
}

// not tested
func getShortcutTarget(path string) string {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}

	var header ShellLinkHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	target := fmt.Sprint(header.FileAttr)
	return target
}

// not tested
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

// not tested
func getShortcuts() map[string]string {
	var shortcuts map[string]string = make(map[string]string)

	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var ShortcutFolder = fmt.Sprintf("%s/Shortcuts", executablePath)

	for _, s := range find(ShortcutFolder, ".lnk") {
		var target = getShortcutTarget(s)
		shortcuts[s] = target
	}

	return shortcuts
}

// not tested
func getAliases() map[string]string {
	var executablePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	file, err := os.Open(fmt.Sprintf("%s/aliases.txt", executablePath))
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
