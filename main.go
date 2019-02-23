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
	pretty = flag.Bool("p", false, "pretty-prints")
)

func doObject(args []string) error {
	jsons := make(map[string]interface{}, len(args))

	for _, arg := range flag.Args() {
		kv := strings.SplitN(arg, "=", 2)
		s := ""
		if len(kv) > 0 {
			s = kv[0]
		}
		if len(kv) != 2 {
			return fmt.Errorf("Argument %q is neither k=v nor k@v", s)
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
		enc := json.NewEncoder(os.Stdout)
		if *pretty {
			enc.SetIndent("", "    ")
		}
		err := enc.Encode(jsons)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	err := doObject(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
