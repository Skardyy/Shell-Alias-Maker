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
## Config  
sam will create a config folder in ~/.sam
you can sam -amend to apply the changed sam config file to your shell file, or sam -add to add to both.  
sam config file is the middle between the cli tool and the shell config file  
### Guidelines  
* all apps inside the ~/.sam will be added automatically to the shell config file (names will be inherited)  
* you can create aliases to them or to other commands by doing:  
```diff
fx : firefox
alias : original_name
fe : fzf --preview "bat --color=always --theme=Dracula {}"
```  
in the above you can see aliases can be created to apps inside the .sam folder or even to full commands used in the shell  
* the [] in the start of the file contains the path to the shellConfigFile.  
you can change it by simply changing it in file or doing the sam -init again and giving it the path arg (sam -init $profile)  

## Cross-platform / shell compatiblity  
in order to support new shells, all need to be done is to creat a new struct that implments the ShellParser interface, then in Shell -> getDynShellParser, change to the desired shell parser. switch between goos.os is possible as well for different shells according to need.