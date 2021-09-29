package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/upamune/jigsaw/drawer"
)

var (
	Version  string
	Revision string

	configPath = flag.String("config", "config.yaml", "path of config")
	noResponse = flag.Bool("no-response", false, "whether to draw response sequences")
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %s\nVersion: %s\nRevision: %s\n", err.Error(), Version, Revision)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	c, err := readConfig(*configPath)
	if err != nil {
		return err
	}

	r, err := readTraceJSON()
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

	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}

	return nil
}

func readTraceJSON() (io.ReadCloser, error) {
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

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
