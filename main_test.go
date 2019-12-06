package main

import (
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
	tests := []struct {
		args  map[string]string
		input []string
		want  int
	}{
		{args: map[string]string{"a": "true", "p": "flase", "v": "false"}, input: []string{`gorilla`, `dog`}, want: 0},
		{args: map[string]string{"p": "true", "a": "flase", "v": "false"}, input: []string{`name=gorilla`}, want: 0},
		{args: map[string]string{"p": "true", "a": "true", "v": "false"}, input: []string{`gorilla`, `cat`}, want: 0},
		{args: map[string]string{"v": "true"}, input: []string{""}, want: 0},
		{args: map[string]string{"p": "false", "a": "false", "v": "false"}, input: []string{`gorilla`}, want: 1},
	}

	var bufstdout, bufstderr bytes.Buffer
	oldstdout, oldstderr := stdout, stderr
	stdout, stderr = &bufstdout, &bufstderr
	defer func() { stdout, stderr = oldstdout, oldstderr }()

	for _, test := range tests {
		os.Args = append([]string{""}, test.input...)

		for arg, value := range test.args {
			flag.CommandLine.Set(arg, value)
		}

		got := run()
		if got != test.want {
			t.Fatalf("want %v, but get %v", test.want, got)
		}
	}
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
