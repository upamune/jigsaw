package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3"
	"github.com/upamune/jigsaw/drawer"
)

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (cli *CLI) Run(args []string) error {
	fs := flag.NewFlagSet("jigsaw", flag.ContinueOnError)

	var (
		configPath = fs.String("config", "", "path of config")
		debug      = fs.Bool("debug", false, "log debug information")

		outputType     = fs.String("type", "mermaid", "output type ('mermaid' or 'plantuml')")
		noResponse     = fs.Bool("no-response", false, "whether to draw response sequences")
		isSkipSelfCall = fs.Bool("skip-self-call", true, "whether to skip self call")
	)

	if err := ff.Parse(fs, args,
		ff.WithEnvVarPrefix("JIGSAW"),
		ff.WithIgnoreUndefined(true),
	); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	c, err := readConfig(*configPath)
	if err != nil {
		return err
	}
	c = overrideConfigWithFlags(c, debug, outputType, noResponse, isSkipSelfCall)

	var filename string
	if args := fs.Args(); len(args) > 0 {
		filename = args[0]
	}

	r, err := readTraceJSON(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	spans, err := parseSpans(r)
	if err != nil {
		return err
	}
	resolvedSpans := resolveSpans(spans)

	var d drawer.Drawer
	switch strings.ToLower(c.Type) {
	case drawer.TypePlantUML, "":
		d = &drawer.PlantUML{}
	case drawer.TypeMermaid:
		d = &drawer.Mermaid{}
	default:
		return fmt.Errorf("unknown type: %s", c.Type)
	}

	s, err := buildDiagram(c, d, resolvedSpans)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(cli.Stdout, s); err != nil {
		return err
	}

	return nil
}

func overrideConfigWithFlags(
	c config,
	debug *bool,
	outputType *string,
	noResponse, isSkipSelfCall *bool,
) config {
	if debug != nil {
		c.Debug = *debug
	}
	if outputType != nil {
		c.Type = *outputType
	}
	if noResponse != nil {
		c.NoResponse = *noResponse
	}
	if isSkipSelfCall != nil {
		c.IsSkipSelfCall = *isSkipSelfCall
	}
	return c
}

func readTraceJSON(filename string) (io.ReadCloser, error) {
	var r io.ReadCloser
	switch filename {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		r = f
	}

	return r, nil
}
