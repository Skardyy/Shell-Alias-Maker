## CLI goals  
```diff
+ open apps by having a folder of shortcuts to those apps
+ have a file called config.txt to alias to those apps
+ have aliases for terminal commands
+ let normal terminal commands through
+ add a way to spin different shells
- add a way to read all installed apps and use them as well (maybe even use them only if its good)  
- add a reload cmd like ?  
- make a way to fetch init commands that install things  
```
## Usage  
* download cc.exe from [latest release](https://github.com/Skardyy/cc/releases/latest).  
Or~  
```diff
git clone https://github.com/Skardyy/cc
cd cc
go build -ldflags "-s -w"
```
* add the dir that contains cc.exe into your path env variable.  
* open your terminal, write cc, and you're done.  
## Config  
create a folder called Apps in the root dir of cc.exe.  
inside that folder you can put your .lnk files, .url files, and any file that requires a simple '. path\to\file' to run  
inside that folder you can also put a config.txt file to create aliases and change the default shell.  
### Config file  
check [Config guidelines](https://github.com/Skardyy/cc/blob/main/Apps/README.md) for detail of how to write a config file.  
### Important  
the Apps folder that contains the apps and config.txt file must be at the same dir as cc.exe  
## to-know  
because cc is an unsigned executable and runs commands, if the command or app you're running is unknown cc will be flagged as antivirus, in such case u can exclude it from the scan of antivirus (or just don't run shady apps / commands), but do so in cautions
