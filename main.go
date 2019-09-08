package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
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
	array   = flag.Bool("a", false, "creates an array of words")
	pretty  = flag.Bool("p", false, "pretty-prints")
	version = flag.Bool("v", false, "show version")

	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

func isRawString(s string) bool {
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return true
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return true
	}
	return false
}

func parseValue(s string) interface{} {
	if s == "" {
		return ""
	}
	if isRawString(s) {
		return json.RawMessage(s)
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	return s
}

func readFile(fname string) (interface{}, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func doArray(args []string) (interface{}, error) {
	jsons := []interface{}{}
	for _, value := range args {
		jsons = append(jsons, parseValue(value))
	}
	return jsons, nil
}

func isKeyFile(s string) bool {
	pos := strings.IndexRune(s, ':')
	return pos > 0 && pos == len(s)-1
}

func doObject(args []string) (interface{}, error) {
	jsons := make(map[string]interface{}, len(args))
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		s := ""
		if len(kv) > 0 {
			s = kv[0]
		}
		if len(kv) != 2 {
			return nil, fmt.Errorf("Argument %q is not k=v", s)
		}
		if isKeyFile(kv[0]) {
			// For argument a:=b, read value from file "b".
			v, err := readFile(kv[1])
			if err != nil {
				return v, err
			}
			key := kv[0][:len(kv[0])-1]
			jsons[key] = v
		} else {
			jsons[kv[0]] = parseValue(kv[1])
		}
	}

	return jsons, nil
}

func doVersion() error {
	enc := json.NewEncoder(stdout)
	if *pretty {
		enc.SetIndent("", "    ")
	}
	return enc.Encode(&Version{
		Program:     "gjo",
		Description: "This is inspired by jpmens/jo",
		Author:      "gorilla0513",
		Repo:        "https://github.com/skanehira/gjo",
		Version:     "1.0.2",
	})
}

func run() int {
	flag.Parse()
	if *version {
		err := doVersion()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}

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
		fmt.Fprintln(stderr, err)
		return 1
	}
	enc := json.NewEncoder(stdout)
	if *pretty {
		enc.SetIndent("", "    ")
	}
	err = enc.Encode(value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
