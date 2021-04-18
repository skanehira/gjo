package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestObject(t *testing.T) {
	tests := []struct {
		input []string
		want  string
		err   string
	}{
		{input: []string{``}, want: ``, err: `Argument "" is not k=v`},
		{input: []string{`a`}, want: ``, err: `Argument "a" is not k=v`},
		{input: []string{`a=`}, want: `{"a":""}`, err: ``},
		{input: []string{`a=1`}, want: `{"a":1}`, err: ``},
		{input: []string{`a=1.1`}, want: `{"a":1.1}`, err: ``},
		{input: []string{`a=true`}, want: `{"a":true}`, err: ``},
		{input: []string{`a=false`}, want: `{"a":false}`, err: ``},
		{input: []string{`a=null`}, want: `{"a":null}`, err: ``},
		{input: []string{`a=s`}, want: `{"a":"s"}`, err: ``},
		{input: []string{`a={"a":"s"}`}, want: `{"a":{"a":"s"}}`, err: ``},
		{input: []string{`a=["a","s"]`}, want: `{"a":["a","s"]}`, err: ``},
		{input: []string{`a:=testdata/test.json`}, want: `{"a":{"foo":true}}`, err: ``},
		{input: []string{`a="+123456"`}, want: `{"a":"+123456"}`, err: ``},
	}

	for _, test := range tests {
		value, err := doObject(test.input)
		if err != nil {
			if err.Error() != test.err {
				t.Fatal(err)
			}
		} else {
			var buf bytes.Buffer
			err = json.NewEncoder(&buf).Encode(value)
			if err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSpace(buf.String())
			if got != test.want {
				t.Fatalf("want %q, but got %q", test.want, got)
			}
		}
	}
}

func TestArray(t *testing.T) {
	tests := []struct {
		input []string
		want  string
		err   string
	}{
		{input: []string{``}, want: `[""]`, err: ``},
		{input: []string{`a`}, want: `["a"]`, err: ``},
		{input: []string{`1`}, want: `[1]`, err: ``},
		{input: []string{`1.1`}, want: `[1.1]`, err: ``},
		{input: []string{`true`}, want: `[true]`, err: ``},
		{input: []string{`false`}, want: `[false]`, err: ``},
		{input: []string{`false`, `1`, `a`}, want: `[false,1,"a"]`, err: ``},
	}

	for _, test := range tests {
		value, err := doArray(test.input)
		if err != nil {
			if err.Error() != test.err {
				t.Fatal(err)
			}
		} else {
			var buf bytes.Buffer
			err = json.NewEncoder(&buf).Encode(value)
			if err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSpace(buf.String())
			if got != test.want {
				t.Fatalf("want %q, but got %q", test.want, got)
			}
		}
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		pretty bool
		want   error
	}{
		{pretty: false, want: nil},
		{pretty: true, want: nil},
	}

	var bufstdout, bufstderr bytes.Buffer
	oldstdout, oldstderr := stdout, stderr
	stdout, stderr = &bufstdout, &bufstderr
	defer func() { stdout, stderr = oldstdout, oldstderr }()

	for _, test := range tests {
		*pretty = test.pretty
		got := doVersion()
		if got != test.want {
			t.Fatalf("want %q, but get %q", test.want, got)
		}
	}

}

func TestRun(t *testing.T) {
	// Take advantage of struct field's zero values, only include fields
	// with non-zero values, be careful about stdin field and ptr values
	// below...
	tests := []struct {
		flags       map[string]string // simulate command line flags
		args        []string          // simulate command line arguments (b=true)
		stdin       strings.Reader    // provide "standard input"
		want_status int               // expected status
		want_stdout string            // expected stdout
	}{
		{flags: map[string]string{"a": "true"}, args: []string{`gorilla`, `dog`}, want_stdout: `["gorilla","dog"]`},
		{flags: map[string]string{"p": "true"}, args: []string{`name=gorilla`},
			want_stdout: `{
    "name": "gorilla"
}`},
		{flags: map[string]string{"p": "true", "a": "true"}, args: []string{`gorilla`, `cat`},
			want_stdout: `[
    "gorilla",
    "cat"
]`},
		{flags: map[string]string{"p": "true", "v": "true"}, args: []string{""},
			want_stdout: `{
    "program": "gjo",
    "description": "This is inspired by jpmens/jo",
    "author": "skanehira",
    "repo": "https://github.com/skanehira/gjo",
    "version": "1.0.3"
}`},
		{args: []string{`gorilla`}, want_status: 1},
		// test that -n works (true -> do not add keys with empty values)
		{args: []string{`name=`}, want_stdout: `{"name":""}`},
		{flags: map[string]string{"n": "true"}, args: []string{`name=`}, want_stdout: `{}`},
		// test that -B works, booleans are treated as strings
		{args: []string{`name=true`}, want_stdout: `{"name":true}`},
		{args: []string{`name=false`}, want_stdout: `{"name":false}`},
		{args: []string{`name=flase`}, want_stdout: `{"name":"flase"}`},
		{args: []string{`name=null`}, want_stdout: `{"name":null}`},
		{flags: map[string]string{"B": "true"}, args: []string{`name=true`}, want_stdout: `{"name":"true"}`},
		{flags: map[string]string{"B": "true"}, args: []string{`name=false`}, want_stdout: `{"name":"false"}`},
		{flags: map[string]string{"B": "true"}, args: []string{`name=null`}, want_stdout: `{"name":"null"}`},
		// Test reads from stdin, with and without -e
		{want_status: 2},
		{flags: map[string]string{"e": "true"}},
		{flags: map[string]string{"a": "true"}, stdin: *strings.NewReader("ape foo"), want_stdout: `["ape","foo"]`},
		{stdin: *strings.NewReader("ape=foo"), want_stdout: `{"ape":"foo"}`},
	}

	var bufstdout, bufstderr bytes.Buffer
	oldstdout, oldstderr := stdout, stderr
	stdout, stderr = &bufstdout, &bufstderr
	defer func() { stdout, stderr = oldstdout, oldstderr }()

	// flags.Usage writes to flag.CommandLine.Output, so capture it
	var cmdoutput bytes.Buffer
	flag.CommandLine.SetOutput(bufio.NewWriter(&cmdoutput))

	// test cases count from n=0...
	for n, test := range tests {
		bufstdout.Reset()
		cmdoutput.Reset()

		os.Args = append([]string{""}, test.args...)
		for test_flag, value := range defaultFlags(test.flags) {
			flag.CommandLine.Set(test_flag, value)
		}

		got_status := run(&test.stdin)
		got_stdout := bufstdout.String()
		got_stdout = strings.TrimSuffix(got_stdout, "\n")
		if got_status != test.want_status {
			t.Fatalf("test_case: %d, want %v, but got %v", n, test.want_status, got_status)
		}
		if got_stdout != test.want_stdout {
			t.Fatalf("test_case: %d, want '%v', but got '%v'", n, test.want_stdout, got_stdout)
		}
	}
}

// defaultFlags merges a set of override values into the standard set of
// command line flags, overriding the default values.
func defaultFlags(overrides map[string]string) map[string]string {
	// HEADS UP, keep this up to date w.r.t. flag settings in main.go
	args := map[string]string{
		"a": "false",
		"B": "false",
		"e": "false",
		"n": "false",
		"p": "false",
		"v": "false",
	}
	for k, v := range overrides {
		args[k] = v
	}
	return args
}

func TestIsKeyFile(t *testing.T) {
	tests := []struct {
		arg  string
		want bool
	}{
		{arg: "", want: false},
		{arg: ":", want: false},
		{arg: "a", want: false},
		{arg: "a:", want: true},
		{arg: "a:=", want: false},
		{arg: "a::", want: false},
	}

	for _, test := range tests {
		got := isKeyFile(test.arg)
		if got != test.want {
			t.Fatalf("want %v, but get %v", test.want, got)
		}
	}
}
