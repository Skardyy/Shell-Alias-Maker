# CC  
cc is used for easily creating perm aliases for terminal
## Usage  
* download cc.exe from [latest release](https://github.com/Skardyy/cc/releases/latest).  
Or ~  
```diff
git clone https://github.com/Skardyy/cc
cd cc
go build
```
* start by doing ./cc.exe -help to see the commands
## Config  
cc will create a config folder in ~/.cc  
you can cc -amend to apply the changed cc config file to your shell file, or cc -add to add to both.  
cc config file is the middle between the cli tool and the shell config file  
### Guidelines  
* all apps inside the ~/.cc will be added automatically to the shell config file (names will be inherited)  
* you can create aliases to them or to other commands by doing:  
```diff
fx : firefox
alias : original_name
fe : fzf --preview "bat --color=always --theme=Dracula {}"
```  
in the above you can see aliases can be created to apps inside the .cc folder or even to full commands used in the shell  
* the [] in the start of the file contains the path to the shellConfigFile.  
you can change it by simply changing it in file or doing the cc -init again and giving it the path arg (cc -init $profile)  
* for other commands simply do cc -h to see a quick description of them.  

## Cross-platform / shell compatiblity  
in order to support new shells, all need to be done is to creat a new struct that implments the ShellParser interface, then in mFile -> populateShellParser, change to the desired shell parser. switch between goos.os is possible as well for different shells according to need.