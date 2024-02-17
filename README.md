# Shell Alias Maker  
sam is used for easily creating perm aliases for terminal
## Usage  
* download sam.exe from [latest release](https://github.com/Skardyy/Shell-Alias-Maker/releases/latest).  
Or ~  
```diff
git clone https://github.com/Skardyy/Shell-Alias-Maker
cd Shell-Alias-Maker
go build
```
* start by doing ./sam.exe -h to see the commands
* you must do sam -init before doing amend (sam must know where your shell config file is located)
## Config  
sam will create a config folder in ~/.sam
you can sam -amend to apply the changed sam config file to your shell file.  
### Guidelines  
* all apps inside the ~/.sam will be added automatically to the shell config file (names will be inherited)  
* you can create aliases to them or to other commands by doing:  
```diff
sam add -alias fx firefox
sam add -alias alias original_name
sam add -alias fzf '--preview "bat --color=always --theme=Dracula {}"'
```  
finally do sam -amend to apply the changes to the shell file. sam -amend also applies manually removed / added files to the ~/.sam folder.
* the same can be done with paths, only that in paths case it won't add it to the shell config file, but just will be there for u to remmember where is ur config files by just doing sam get -path

## Cross-platform / shell compatiblity  
in order to support new shells, all need to be done is to creat a new struct that implments the ShellParser interface, then in Shell -> getDynShellParser, change to the desired shell parser. switch between goos.os is possible as well for different shells according to need.