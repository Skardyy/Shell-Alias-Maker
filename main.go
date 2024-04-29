package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

var aCmd *flag.FlagSet
var cCmd *flag.FlagSet
var gCmd *flag.FlagSet
var rCmd *flag.FlagSet
var cf configFile
var apps, _ = getApps()

func main() {
	cf = configFile{}
	err := cf.readConfig()
	if err != nil {
		fmt.Println("had a problem parsing the config file: ", err)
	}

	cCmd = flag.NewFlagSet("create", flag.ExitOnError)
	cBaseDir := cCmd.String("baseDir", "~/desktop", "the base dir to walk from")
	cSuffix := cCmd.String("suffix", ".url .exe .lnk", "all the suffixes to capture")
	cRecursive := cCmd.Bool("recursive", false, "should the app walk new dir found on baseDir")

	aCmd = flag.NewFlagSet("add", flag.ExitOnError)
	aPath := aCmd.Bool("path", false, "adds a path to the config, followed by: key:string target:string")
	aAlias := aCmd.Bool("alias", false, "adds a alias to the config, followed by: key:string target:string")

	rCmd = flag.NewFlagSet("rm", flag.ExitOnError)
	rPath := rCmd.Bool("path", false, "removes a path from the config, followed by: key:string")
	rAlias := rCmd.Bool("alias", false, "removes a alias from the config, followed by: key:string")

	gCmd = flag.NewFlagSet("get", flag.ExitOnError)
	gDir := gCmd.Bool("dir", false, "gets the dir of all sam configs")
	gAlias := gCmd.Bool("alias", false, "gets all the created aliases")
	gPaths := gCmd.Bool("path", false, "gets all stored paths")
	gApps := gCmd.Bool("app", false, "gets all stored apps")

	flag.BoolFunc("clear", "removes all the content sam created in the shell config file", clear)
	flag.Func("init", "creates (if dosen't exists) a ~/.sam folder and inside of it config.txt file", Init)
	flag.BoolFunc("amend", "amends manually deleted/added apps to the config file\nthen amends the changes made in config.json to the shellConfig file\n", amend)

	flag.Usage = func() {
		fmt.Println("Global funcs:")
		flag.PrintDefaults()
		printBr()

		fmt.Println("add:")
		aCmd.PrintDefaults()
		printBr()

		fmt.Println("rm:")
		rCmd.PrintDefaults()
		printBr()

		fmt.Println("create:")
		cCmd.PrintDefaults()
		printBr()

		fmt.Println("get:")
		gCmd.PrintDefaults()
		printBr()
	}
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "get":
		gCmd.Parse(os.Args[2:])
		handleGet(gDir, gAlias, gPaths, gApps)
	case "rm":
		rCmd.Parse(os.Args[2:])
		handleRemove(rPath, rAlias, rCmd.Args())
	case "create":
		cCmd.Parse(os.Args[2:])
		handleCreate(cBaseDir, cSuffix, cRecursive)
	case "add":
		aCmd.Parse(os.Args[2:])
		handleAdd(aPath, aAlias, aCmd.Args())
	case "-clear":
	case "-init":
	case "-amend":
	default:
		if value, flag := cf.Paths[os.Args[1]]; flag {
			fmt.Println(value)
			os.Exit(0)
		}
		fmt.Println(os.Args[1], "Isn't a command, try sam -h")
		os.Exit(0)
	}
}

func printBr() {
	fd := int(os.Stdout.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		fmt.Println(err)
		width = 100
	}

	line := strings.Repeat("-", width)
	fmt.Println(line)
}

func handleAdd(path *bool, alias *bool, args []string) error {
	if len(args) < 2 || (!*path && !*alias) {
		aCmd.PrintDefaults()
		os.Exit(0)
	}
	key, target := args[0], args[1]
	if *path {
		if _, flag := cf.Paths[key]; flag {
			var f string
			fmt.Println(key, "already exists in your stored paths, override it? y/n")
			fmt.Scanln(&f)
			if f != "y" {
				os.Exit(0)
			}
		}
		cf.Paths[key] = target
		err := cf.writeConfig()
		if err != nil {
			return err
		}
	}
	if *alias {
		if _, flag := cf.Aliases[key]; flag {
			var f string
			fmt.Println(key, "already exists in your stored aliases, override it? y/n")
			fmt.Scanln(&f)
			if f != "y" {
				os.Exit(0)
			}
		}
		cf.Aliases[key] = target
		err := cf.writeConfig()
		if err != nil {
			return err
		}
		fmt.Println("please ammend to change the shell")
	}
	return nil
}

func handleRemove(path *bool, alias *bool, args []string) error {
	if len(args) < 1 || (!*path && !*alias) {
		rCmd.PrintDefaults()
		os.Exit(0)
	}
	key := args[0]
	if *path {
		delete(cf.Paths, key)
	}
	if *alias {
		delete(cf.Aliases, key)
	}
	err := cf.writeConfig()
	if err != nil {
		return err
	}
	fmt.Println("please ammend to change the shell")
	return nil
}

func handleCreate(baseDir *string, suffix *string, recursive *bool) error {
	suffixes := strings.Split(*suffix, " ")
	aliases, err := walkBaseDir(*baseDir, suffixes, *recursive)
	if err != nil {
		return err
	}

	for key, target := range aliases {
		if _, flag := cf.Apps[key]; !flag {
			newPath, err := storePath(target)
			if err != nil {
				return err
			}
			base := filepath.Base(newPath)
			newKey := strings.TrimSuffix(base, filepath.Ext(base))
			cf.Apps[newKey] = newPath
		} else {
			fmt.Println(key, "already exists, hence overriding it")
		}
	}

	err = cf.writeConfig()
	if err != nil {
		return err
	}
	fmt.Println("Please use the amend command to apply the changes")
	return nil
}

func handleGet(dir *bool, alias *bool, paths *bool, apps *bool) error {
	if !*dir && !*alias && !*paths && !*apps {
		gCmd.PrintDefaults()
	}
	if *dir {
		path, err := getConfigDirPath()
		if err != nil {
			return err
		}
		fmt.Println(path)
	}
	if *alias {
		fmt.Println(cf.echo(cf.Aliases))
	}
	if *paths {
		fmt.Println(cf.echo(cf.Paths))
	}
	if *apps {
		fmt.Println(cf.echo(cf.Apps))
	}
	return nil
}

func Init(newShellConfigPath string) error {
	cf.ShellConfigPath = newShellConfigPath
	err := cf.writeConfig()

	return err
}
func amend(s string) error {
	err := cf.ammend(apps)
	if err != nil {
		return err
	}
	parser := populateShellParser(cf)
	err = parser.confirm()
	if err != nil {
		return err
	}
	err = createReproduceFile(parser.reproduceContent)
	if err != nil {
		return err
	}
	fmt.Println("successfully amended, please rerun your shell config file")
	return nil
}
func clear(s string) error {
	parser := getDynShellParser(cf)
	return parser.confirm()
}
