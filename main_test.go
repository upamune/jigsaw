package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/tenntenn/golden"
)

var flagUpdate bool

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestCLI_Run(t *testing.T) {
	t.Parallel()

	for _, outputType := range []string{"mermaid", "plantuml"} {
		outputType := outputType
		t.Run(outputType, func(t *testing.T) {
			t.Parallel()

			cases := map[string]struct {
				args    string
				in      string
				wantErr bool
			}{
				"simple": {
					args: "-config ./testdata/simple/config.yaml ./testdata/simple/trace.json",
				},
				"simple:no-response": {
					args: "-config ./testdata/simple/config.yaml -no-response ./testdata/simple/trace.json",
				},
				"simple:no-response-not-skip-self-call": {
					args: "-config ./testdata/simple/config.yaml -no-response -skip-self-call=false ./testdata/simple/trace.json",
				},
				"simple:not-no-response-not-skip-self-call": {
					args: "-config ./testdata/simple/config.yaml -no-response=false -skip-self-call=false ./testdata/simple/trace.json",
				},
				"complex": {
					args: "./testdata/complex/trace.json",
				},
				"complex:no-response": {
					args: "-no-response ./testdata/complex/trace.json",
				},
				"complex:no-response-not-skip-self-call": {
					args: "-no-response -skip-self-call=false ./testdata/complex/trace.json",
				},
				"complex:not-no-response-not-skip-self-call": {
					args: "-no-response=false -skip-self-call=false ./testdata/complex/trace.json",
				},
			}

			for name, tt := range cases {
				name, tt := name, tt
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					var stdout, stderr bytes.Buffer
					cli := &CLI{
						Stdout: &stdout,
						Stderr: &stderr,
						Stdin:  strings.NewReader(tt.in),
					}

					args := append([]string{"-type", outputType}, strings.Split(tt.args, " ")...)
					err := cli.Run(args)

					if tt.wantErr {
						if err == nil {
							t.Error("want error, but got nil")
						}
					} else {
						if err != nil {
							t.Errorf("want nil, but got an error: %v", err)
						}
					}

					name = fmt.Sprintf("%s_%s", outputType, name)
					c := golden.New(t, flagUpdate, "testdata/golden", name)

					if diff := c.Check("_stdout", &stdout); diff != "" {
						t.Error("stdout\n", diff)
					}

					if diff := c.Check("_stderr", &stderr); diff != "" {
						t.Error("stderr\n", diff)
					}
				})
			}
		})
	}
}
