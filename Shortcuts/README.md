## Shortcuts  
in order to use the shortcuts future, you need to add .lnk files into this folder.  
the program will detect them on its own,  
names of the shortcuts will be inherited  

## Aliases  
you can write an config.txt file to create aliases to those shortcuts and commands  
### alias to .lnk file  
fx : firefox 
* the above is used to create a shortcut to a .lnk file called firefox  
### alias to terminal command  
ef : fzf | split-path | % { code $_ }  
* the above is used to create a shortcut to the command (with args):> fzf | split-path | % { code $_ }  
### alias to a command + run async  
gn : Google chrome ! async  
* the above is used to create a shortcut to a .lnk file and specifies to run it async  
## Shell  
in the config.txt you can specify the shell to use (defaults to powershell if no shell specifies)  
simply put:  
[<shell>] in the first line of config.txt
