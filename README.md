## CLI goals  
```diff
+ open apps by having a folder of shortcuts to those apps
+ have a file called aliases.txt to alias to those apps
+ have aliases for terminal commands
+ let normal terminal commands through
+ add a way to spin different shells
- create a way to make commands that install dependencies and git projects
! check why fe -a dosent work as a alias and fe-a does.  
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
### Important  
the shortcuts folder that contains the .lnk files and aliases.txt file must be at the same dir as cc.exe  
## to-know  
* because cc is an unsigned executable and runs commands, if the command or app you're running is unknown cc will be flagged as antivirus, in such case u can exclude it from the scan of antivirus (or just don't run shady apps / commands), but do so in cautions
