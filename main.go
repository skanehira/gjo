package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	array              = flag.Bool("a", false, "creates an array of words")
	disable_bool       = flag.Bool("B", false, "disable treating true/false as bool")
	ignore_empty_stdin = flag.Bool("e", false, "empty stdin is not an error")
	ignore_empty_keys  = flag.Bool("n", false, "ignore keys with empty values")
	pretty             = flag.Bool("p", false, "pretty-prints")
	version            = flag.Bool("v", false, "show version")

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
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
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
		if *disable_bool {
			return s
		}
		return true
	}
	if s == "false" {
		if *disable_bool {
			return s
		}
		return false
	}
	if s == "null" {
		if *disable_bool {
			return s
		}
		return nil
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
			value := parseValue(kv[1])
			if value != "" || !*ignore_empty_keys {
				jsons[kv[0]] = value
			}
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
		Author:      "skanehira",
		Repo:        "https://github.com/skanehira/gjo",
		Version:     "1.0.3",
	})
}

func run(stdin io.Reader) int {
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
		dat, err := ioutil.ReadAll(stdin)
		if err != nil {
			fmt.Fprintln(stderr, err)
		}
		args = strings.Fields(string(dat))
		if len(args) == 0 {
			if *ignore_empty_stdin {
				return 0
			}
			flag.Usage()
			return 2
		}
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
	os.Exit(run(os.Stdin))
}
