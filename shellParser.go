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
func (scp *ShellConfigParser) confirm() {
	replaceFilePartition(scp.ShellParser.GetPartitionDel(), scp.partitionedContent, getFile(scp.shellConfigPath))
}

type PwshConfigParsser struct {
}

func (psp *PwshConfigParsser) Add(content []string, alias Alias) []string {
	newValue := "\n" + "New-Alias -Name " + alias.name + " -Value " + alias.target + " #" + alias.name
	content = append(content, newValue)
	return content
}
func (psp *PwshConfigParsser) GetPartitionDel() string {
	return "#CC"
}
