package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

var aCmd *flag.FlagSet
var gCmd *flag.FlagSet

func main() {
	aCmd = flag.NewFlagSet("add", flag.ExitOnError)
	aBaseDir := aCmd.String("baseDir", "~/desktop", "the base dir to walk from")
	aSuffix := aCmd.String("suffix", `".url .exe .lnk"`, "all the suffixes to capture")
	aRecursive := aCmd.Bool("recursive", false, "should the app walk new dir found on baseDir")

	gCmd = flag.NewFlagSet("get", flag.ExitOnError)
	gDir := gCmd.Bool("dir", false, "gets the dir of all cc configs")
	gAlias := gCmd.Bool("alias", false, "gets all the created aliases")

	flag.BoolFunc("clear", "removes all the content cc created in the shell config file", clear)
	flag.Func("init", "creates (if dosen't exists) a ~/.cc folder and inside of it config.txt file", Init)
	flag.BoolFunc("amend", "amends manually deleted/added apps to the config file\nthen amends the changes made in config.txt to the shellConfig file\n", amend)

	flag.Usage = func() {
		fmt.Println("Global funcs:")
		flag.PrintDefaults()

		fmt.Println("<add>:")
		aCmd.PrintDefaults()

		fmt.Println("<get>:")
		gCmd.PrintDefaults()
	}
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("cc -h")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "get":
		gCmd.Parse(os.Args[2:])
		handleGet(gDir, gAlias)
	case "add":
		aCmd.Parse(os.Args[2:])
		handleCreate(aBaseDir, aSuffix, aRecursive)
	case "-clear":
	case "-init":
	case "-amend":
	default:
		fmt.Println(os.Args[1])
		fmt.Println("\ncc -h")
		os.Exit(0)
	}
}

func handleCreate(baseDir *string, suffix *string, recursive *bool) {
	suffixes := strings.Split(*suffix, " ")
	aliases, err := walkBaseDir(*baseDir, suffixes, *recursive)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	for _, a := range aliases {
		newPath, err := storePath(a.target)
		if err != nil {
			panic(err)
		}
		buffer.WriteString(a.name + " : " + newPath + "\n")
	}

	file, err := getConfigFile()
	if err != nil {
		panic(err)
	}

	err = replaceFilePartition("#Apps", file, true, buffer.String())
	if err != nil {
		panic(err)
	}
	fmt.Println("Please use the amend command to apply the changes")
}

func handleGet(dir *bool, alias *bool) {
	if *dir {
		path, err := getConfigDirPath()
		if err != nil {
			panic(err)
		}
		fmt.Println(path)
	}
	if *alias {
		fmt.Println(echoAliases())
	}
}

func Init(newShellConfigPath string) error {
	err := createConfigDir()
	if err != nil {
		return err
	}
	file, err := getConfigFile()
	if err != nil {
		return err
	}
	err = initConfig(newShellConfigPath, file)
	if err != nil {
		return err
	}

	return nil
}
func amend(s string) error {
	err := amendApps()
	if err != nil {
		panic(err)
	}
	parser := populateShellParser()
	return parser.confirm()
}
func clear(s string) error {
	parser := ShellConfigParser{}
	return parser.confirm()
}
