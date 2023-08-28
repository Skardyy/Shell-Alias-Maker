## CLI goals  
* open apps by having a folder of shortcuts to those apps -> done
* have a file called aliases.txt to alias to those apps -> done
* have file explorer -> done
* let normal terminal commands through -> done  
## Usage  
download cc.exe from [latest release](https://github.com/Skardyy/cc/releases/latest).  
once downloaded add the dir that contains cc.exe into your path env variable.  
open your terminal, write cc, and you're done.  
### Important  
the shortcuts folder that contains the .lnk files and aliases.txt file must be at the same dir as cc.exe  
## to-know  
* cc dosent support io related commands.. all commands are ran asynchronous so cli apps most likely will break it.  
* because cc is an unsigned executable and runs commands, if the command you're running is unknown (custom .lnk files that aren't known apps), cc will be flagged as antivirus, in such case u can exclude it from the scan of antivirus, but do so in cautions
