## Apps  
in order to add apps to cc, you need to add .lnk/.url/any file that can run by simply doing '. path\to\file' files into this folder.  
the program will detect them on its own,  
names of the files will be inherited  

## Aliases  
you can write an config.txt file to create aliases to those shortcuts and commands  
### alias to an app  
fx : firefox 
* the above is used to create a shortcut to an app called firefox  
### alias to terminal command  
ef : fzf | split-path | % { code $_ }  
* the above is used to create a shortcut to the command (with args):> fzf | split-path | % { code $_ }  
### alias to a command/app + run async  
gn : Google chrome ! async  
* the above is used to create a shortcut to a app and specifies it to run it async  
## Shell  
in the config.txt you can specify the shell to use (defaults to powershell if no shell specifies)  
simply put:  
[Shell] in the first line of config.txt
