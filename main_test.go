package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestObject(t *testing.T) {
	tests := []struct {
		input []string
		want  string
		err   string
	}{
		{input: []string{``}, want: ``, err: `Argument "" is neither k=v nor k@v`},
		{input: []string{`a`}, want: ``, err: `Argument "a" is neither k=v nor k@v`},
		{input: []string{`a=`}, want: `{"a":null}`, err: ``},
		{input: []string{`a=1`}, want: `{"a":1}`, err: ``},
		{input: []string{`a=1.1`}, want: `{"a":1.1}`, err: ``},
		{input: []string{`a=true`}, want: `{"a":true}`, err: ``},
		{input: []string{`a=false`}, want: `{"a":false}`, err: ``},
		{input: []string{`a=s`}, want: `{"a":"s"}`, err: ``},
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
		{input: []string{``}, want: `[null]`, err: ``},
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
