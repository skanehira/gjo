# gjo

[![GitHub Actions status](https://github.com/skanehira/gjo/workflows/Go/badge.svg)](https://github.com/skanehira/gjo/actions)
[![CircleCI](https://img.shields.io/circleci/project/github/skanehira/gjo/master.svg)](https://circleci.com/gh/skanehira/gjo/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/skanehira/gjo)](https://goreportcard.com/report/github.com/skanehira/gjo)

Small utility to create JSON objects.  
This was inspired by [jpmens/jo](https://github.com/jpmens/jo).

![sreenshot](./screenshot.png)

## Support OS
- Mac
- Linux
- Windows

## Requirements
- Go 1.1.14~
- Git

## Installtion
### Build
```sh
$ git clone https://github.com/skanehira/gjo.git
$ cd gjo
$ GO111MODULE=on go install
```

### Binary
Please download from [releases](https://github.com/skanehira/gjo/releases)

## Usage
### Mac and Linux
```sh
$ gjo -p status=$(gjo name=gorilla age=26 lang=$(gjo -a Go Java PHP))
{
    "status": {
        "age": 26,
        "lang": [
            "Go",
            "Java",
            "PHP"
        ],
        "name": "gorilla"
    }
}
$ gjo -h
Usage of ./gjo:
  -B    disable treating true/false as bool
  -a    creates an array of words
  -e    empty stdin is not an error
  -n    ignore keys with empty values
  -p    pretty-prints
  -v    show version
```

### Windows
If you want to use `$()` on the Windows, please install [shellwrap](https://github.com/mattn/shellwrap).

```sh
shellwrap gjo -p status=$(gjo name=gorilla age=26)
```

## Author
gorilla0513
