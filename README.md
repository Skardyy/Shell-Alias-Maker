# Shell Alias Maker  
Sam is used for easily creating shell scripts that run apps or commands and can be ported everywhere  

## Usage  
* download sam.exe from [latest release](https://github.com/Skardyy/Shell-Alias-Maker/releases/latest).  
Or ~  
```pwsh
git clone https://github.com/Skardyy/Shell-Alias-Maker
cd Shell-Alias-Maker
go build
```
* add the directory that contains sam.exe file into your env variables
* start by doing `sam -h` to see the commands

## Config  
sam will create a dir in ~/.sam where it going to store the config and its prequisites  
it will include a pre and dst dir, `the dst dir is where the scripts will be generated so you should add this dir into your env path`  
the pre dir is sam will hold file that were copied over like .lnk .url .exe files and later it going to point to them for easier configuration  
the pre dir is also where you can add files manually and when doing `sam -a ps1` sam will account for all the manual changes that happend there and in the config.json file  
> \[!Note]  
> as mentioned above any manual changes made to the pre folder or the config.json file will account after doing sam -a  
> yet you will be prompted when there is an conflict between an existing shell scripts and what is specified in the config.json

## The Idea  
the idea behind this app was to create a portable way to creating shell scripts to run apps and commands like aliases  
the commands should be simple enough that every shell no matter what platform you use should understand them and act the same  
some shells tho (nushell for instance) don't support running their own shell scripts by using only name so you may need to search for another tool in those shells  
