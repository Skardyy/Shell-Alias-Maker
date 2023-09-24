package main

type Alias struct {
	name   string
	target string
}

type ShellConfigParser struct {
	shellConfigPath    string
	partitionedContent []string
	ShellParser
}

type ShellParser interface {
	Add(content []string, alias Alias) []string
	GetPartitionDel() string
}

func (scp *ShellConfigParser) With(shellConfigPath string, shellParser ShellParser) {
	scp.shellConfigPath = shellConfigPath
	scp.ShellParser = shellParser
}
func (scp *ShellConfigParser) RemoveAll() {
	scp.partitionedContent = nil
}
func (scp *ShellConfigParser) Add(aliases ...Alias) {
	for _, a := range aliases {
		scp.partitionedContent = scp.ShellParser.Add(scp.partitionedContent, a)
	}
}
func (scp *ShellConfigParser) confirm() error {
	file, err := getFile(scp.shellConfigPath)
	if err != nil {
		return err
	}
	err = replaceFilePartition(scp.ShellParser.GetPartitionDel(), file, false, scp.partitionedContent...)
	if err != nil {
		return err
	}
	return nil
}

func getDynShellParser() ShellConfigParser {
	//---------- here give other code different parsers ----------
	// can switch between goos.os for different ones
	parser := ShellConfigParser{}
	parser.With(shellConfigPath, &PwshConfigParsser{})
	return parser
}

// -------------------- create new config parsers to support more shells --------------------

// ---------- PowerShell file parser ----------
type PwshConfigParsser struct {
}

func (psp *PwshConfigParsser) Add(content []string, alias Alias) []string {
	newValue := "\nfunction " + alias.name + " { " + alias.target + " }"
	content = append(content, newValue)
	return content
}
func (psp *PwshConfigParsser) GetPartitionDel() string {
	return "#SAM"
}

// ---------- PowerShell file parser ----------
