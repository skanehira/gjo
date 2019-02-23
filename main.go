package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

// Version gjo version info
type Version struct {
	Program     string `json:"program"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Repo        string `json:"repo"`
	Version     string `json:"version"`
}

var (
	format  = flag.Bool("f", false, "format json")
	version = flag.Bool("v", false, "show version")
)

func main() {
	flag.Parse()
	args := flag.Args()

	// if version flag is true, display version info
	if *version {
		b, err := json.Marshal(&Version{
			Program:     "gjo",
			Description: "This is inspired by jpmens/jo",
			Author:      "gorilla0513",
			Repo:        "https://github.com/skanehira/gjo",
			Version:     "1.0.0",
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
		return
	}

	// parse args to map
	jsons := make(map[string]interface{}, len(args))

	for _, arg := range flag.Args() {
		kv := strings.Split(arg, "=")
		if len(kv) == 2 {
			key, value := kv[0], kv[1]

			switch value {
			case "true":
				jsons[key] = true
			case "false":
				jsons[key] = false
			default:
				jsons[key] = value
			}
		}
	}

	if len(jsons) != 0 {
		// parse map to json string
		var output string
		j, err := json.Marshal(jsons)

		if err != nil {
			panic(err)
		}

		// if format flag is true, format json string
		if *format {
			out := new(bytes.Buffer)
			json.Indent(out, j, "", "    ")
			output = out.String()
		} else {
			output = string(j)
		}

		fmt.Println(output)
	}
}
