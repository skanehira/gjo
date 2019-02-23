# gjo
Small utility to create JSON objects.  
This was inspired by [jpmens/jo](https://github.com/jpmens/jo).

![sreenshot](./screenshot.png)

## Requirements
- Go 1.1.14~
- Git

## Installtion
```sh
$ git clone https://github.com/skanehira/gjo.git
$ cd gjo
$ GO111MODULE=on go install
```

## Usage
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
Usage of gjo:
  -a    creates an array of words
  -p    pretty-prints
  -v    show version
```

## Author
gorilla0513
