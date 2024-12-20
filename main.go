package main

import (
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/fatih/color"
)

var cf configFile
var apps map[string]string
var captureUsuage func()
var addUsuage func()

func main() {
  cf = configFile{}
  err := cf.readConfig()
  apps, _ = getApps()
  if err != nil {
    fmt.Println("had a problem parsing the config file: ", err)
  }

  captureUsuage = func() {
    fmt.Println("capture:")
    fmt.Println("  -dir          the base dir to search from")
    fmt.Println("  -ext          (optional) the suffixes to search for, default is: '.lnk .url .exe'")
    fmt.Println("  -r            (optional) rather or not to walk the path recursively")
    fmt.Println("  -c            (optional) rather or not to copy the captured file")
  }
  addUsuage = func() {
    fmt.Println("add:")
    fmt.Println("  -n         the name of shell script")
    fmt.Println("  -t       the content of the shell script (either a script or a path to an executable)")
    fmt.Println("  -c            (optional) rather or not to copy the captured file")
  }
  usuage := func() {
    captureUsuage()

    fmt.Println("")
    addUsuage()

    fmt.Println("")
    fmt.Println("-a       applies the config.json to create shell scripts, requires file extension (like .nu / .bat / .ps1 et..)")
    fmt.Println("-p       prints the config.json content formatted")
  }

  if len(os.Args) < 2 {
    usuage()
    os.Exit(0)
  }

  switch os.Args[1] {
  case "-h":
    usuage()
    os.Exit(0)
  case "capture":
    capture()
    os.Exit(0)
  case "add":
    add()
    os.Exit(0)
  case "a":
    apply()
    os.Exit(0)
  case "-print":
    printConfig()
    os.Exit(0)
  case "print":
    printConfig()
    os.Exit(0)
  case "-p":
    printConfig()
    os.Exit(0)
  case "p":
    printConfig()
    os.Exit(0)
  default:
    fmt.Println(os.Args[1], "Isn't a command, try sam -h")
    os.Exit(0)
  }
}

func findArgAndNext(args []string, targetArg string, onlyArg bool) (string, bool) {
  for i := 0; i < len(args); i++ {
    if args[i] == targetArg && onlyArg {
      return "", true
    }
    if args[i] == targetArg {
      // Check if there is a next argument
      if i+1 < len(args) {
        return args[i+1], true
      } else {
        // Target found but no next argument
        return "", false
      }
    }
  }
  // Target not found
  return "", false
}

func printConfig() error {
  content, err := json.MarshalIndent(cf, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(content))
  return nil
}

func add() error {
  args := os.Args
  name, nameOk := findArgAndNext(args, "-n", false)
  target, targetOk := findArgAndNext(args, "-t", false)
  _, copyOk := findArgAndNext(args, "-c", false)

  if !nameOk || !targetOk {
    addUsuage()
    os.Exit(0)
  }

  //check if already exists in the config.json file
  _, added := addSingle(name, target, copyOk)
  if !added {
    os.Exit(0)
  }
  err := cf.writeConfig()
  if err != nil {
    return err
  }
  color.Magenta("please apply to see the changes")
  return nil
}

func addSingle(name string, path string, copyFlag bool) (string, bool) {
  if _, flag := cf.Commands[name]; flag {
    f := promptTillSuccess("already exists in your stored commands, override it? y/n")
    if f != "y" {
      return "", false
    }
  }

  newPath := path
  if copyFlag {
    attempPath, err := storePath(path)
    if err != nil {
      return "", false
    }
    newPath = attempPath
  }
  //write it to the config.json file
  cf.Commands[name] = newPath
  return newPath, true
}

func capture() error {
  args := os.Args
  dir, dirOk := findArgAndNext(args, "-dir", false)
  ext, extOk := findArgAndNext(args, "-ext", false)
  _, recursiveOk := findArgAndNext(args, "-r", true)
  _, copyOk := findArgAndNext(args, "-c", true)
  suffixes := []string{".url", ".lnk", ".exe"}
  if extOk {
    suffixes = strings.Split(ext, " ")
  }

  if !dirOk {
    captureUsuage()
    os.Exit(0)
  }

  aliases, err := walkBaseDir(dir, suffixes, recursiveOk)
  if err != nil {
    return err
  }

  for key, target := range aliases {
    newPath, added := addSingle(key, target, copyOk)
    if !added {
      continue
    }
    base := filepath.Base(newPath)
    key := strings.TrimSuffix(base, filepath.Ext(base))
    cf.Commands[key] = newPath
  }

  err = cf.writeConfig()
  if err != nil {
    fmt.Println("had an issue writing to config.json", err)
  } else {
    color.Magenta("please apply to see the changes")
  }
  return nil
}

func apply() error {
  ext, extOk := findArgAndNext(os.Args, "-a", false)
  if !extOk {
    ext, extOk = findArgAndNext(os.Args, "a", false)
  }

  if !extOk {
    fmt.Println("apply requires file extension specified, please see sam -h")
    os.Exit(0)
  }

  for key, target := range apps {
    if _, ok := cf.Commands[key]; !ok {
      color.Green("adding " + key + " to the config file")
      cf.Commands[key] = target
    }
  }

  for key, target := range cf.Commands {
    alias := Alias{Name: key, Target: target}
    err := parseAlias(alias, ext, &cf)
    if err != nil {
      fmt.Println("error applying", key, err)
    }
  }

  err := cf.writeConfig()
  if err != nil {
    fmt.Println("had an issue writing to config.json", err)
  }

  return nil
}
