package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

var (
	format = flag.Bool("f", false, "format json")
)

func main() {
	flag.Parse()
	args := flag.Args()

	jsons := make(map[string]interface{}, len(args))

	for _, arg := range flag.Args() {
		kv := strings.SplitN(arg, "=", 2)
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
		var output string
		j, err := json.Marshal(jsons)

		if err != nil {
			panic(err)
		}

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
