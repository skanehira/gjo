package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

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
		output, err := json.Marshal(jsons)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(output))
	}
}
