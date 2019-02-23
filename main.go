package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	array  = flag.Bool("a", false, "creates an array of words")
	pretty = flag.Bool("p", false, "pretty-prints")
)

func doArray(args []string) (interface{}, error) {
	jsons := []interface{}{}

	for _, value := range flag.Args() {
		if value == "" {
			jsons = append(jsons, nil)
			continue
		}
		if value == "true" {
			jsons = append(jsons, true)
			continue
		}
		if value == "false" {
			jsons = append(jsons, false)
			continue
		}

		f, err := strconv.ParseFloat(value, 64)
		if err == nil {
			jsons = append(jsons, f)
			continue
		}
		jsons = append(jsons, value)
	}

	return jsons, nil
}

func doObject(args []string) (interface{}, error) {
	jsons := make(map[string]interface{}, len(args))

	for _, arg := range flag.Args() {
		kv := strings.SplitN(arg, "=", 2)
		s := ""
		if len(kv) > 0 {
			s = kv[0]
		}
		if len(kv) != 2 {
			return nil, fmt.Errorf("Argument %q is neither k=v nor k@v", s)
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

		f, err := strconv.ParseFloat(value, 64)
		if err == nil {
			jsons[key] = f
			continue
		}
		jsons[key] = value
	}

	return jsons, nil
}

func run() int {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 2
	}

	var value interface{}
	var err error

	if *array {
		value, err = doArray(args)
	} else {
		value, err = doObject(args)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	enc := json.NewEncoder(os.Stdout)
	if *pretty {
		enc.SetIndent("", "    ")
	}
	err = enc.Encode(value)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
