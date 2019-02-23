package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
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
		if len(kv) != 2 {
			log.Fatal("Argument `a' is neither k=v nor k@v")
		}
		key, value := kv[0], kv[1]

		if value == "" {
			jsons[key] = nil
			continue
		}
		if value == "true" {
			jsons[key] = true
			continue
		}
		if value == "false" {
			jsons[key] = false
			continue
		}

		f, err := strconv.ParseFloat(key, 64)
		if err == nil {
			jsons[key] = f
			continue
		}
		jsons[key] = value
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
