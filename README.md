

# photoman

A command line tool in Go for managing files. It supports file de-duplication and renaming.

This tool was originally built for managing photo files taken by digital cameras. 

We have > 70k photos taken in the last 10 years stored in a portable hard drive. This tool was made to help better organising
the photos ( sorting by year and month ) and removing duplicated photos due to sharing, coping and multiple backups. 


# Installation

```
go get github.com/dlin-me/photoman

$GOPATH/bin/photoman
```

# Usage

```
NAME:
   Photo Manager - Organise your photo files

USAGE:
   Photo Manager [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   index	Build an index for files in the current directory
   dd		Remove duplicated files
   move		Move files with exif data to proposed directory YYYY/YYYY_MM
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

# Examples

1. Change directory to the hard drive folder, e.g.

    ````
    cd /Volumes/Backup
    ````
    
1. Index files

    ```
    photoman index
    ```
    
1. Deduplicate files in dry mode. All duplicated files will be listed

    ```
    photoman dd -d
    ```
   
1. Deduplicate files

    ```
    photoman dd
    ```
   
1. Move files based on exif date in dry mode. It will list files to move moved. Files will be moved to YYYY/YYYY_MM directories, for example 2016/2016_01

    ```
    photoman move -d
    ```
   
1. Move files based on exif date. Only images files with exif data can be moved

    ```
    photoman move
    ```
   
1. You can also move files based on modified date instead if exif data is not available, however, you don't normally want to move all files in that manner. That's why it requires a directory name to specify the files to be moved.

    ```
    photoman move -g incoming
    ```
