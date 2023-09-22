package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	aCmd := flag.NewFlagSet("add", flag.ExitOnError)
	aBaseDir := aCmd.String("baseDir", "~/desktop", "the base dir to walk from")
	aSuffix := aCmd.String("suffix", `".url .exe .lnk"`, "all the suffixes to capture")
	aRecursive := aCmd.Bool("recursive", false, "should the app walk new dir found on baseDir")

	gCmd := flag.NewFlagSet("get", flag.ExitOnError)
	gDir := gCmd.Bool("dir", false, "gets the dir of all cc configs")
	gAlias := gCmd.Bool("alias", false, "gets all the created aliases")

	flag.BoolFunc("clear", "removes cc config content and removes all the content cc created in the shell config file", clear)
	flag.BoolFunc("init", "creates (if dosen't exists) a ~/.cc folder and inside of it config.txt file", Init)
	flag.BoolFunc("amend", "amends the changes made in config.txt to the shellConfig file", amend)
	flag.BoolFunc("help", "the help command for the cli tool", help)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("cc -help for help")
	}

	switch os.Args[1] {
	case "get":
		handleGet(gDir, gAlias)
	case "add":
		handleCreate(aBaseDir, aSuffix, aRecursive)
	}
}

func handleCreate(baseDir *string, suffix *string, recursive *bool) {

}

func handleGet(dir *bool, alias *bool) {
	if *dir {
		fmt.Println(getConfigDirPath())
	}
	if *alias {
		fmt.Println(echoAliases())
	}
}

func Init(s string) error {
	createConfigDir()

}
func amend(s string) error {

}
func help(s string) error {
	flag.PrintDefaults()
	return nil
}
func clear(s string) error {

}
