package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"

	"github.com/goccy/go-yaml"
	"github.com/mattn/go-jsonpointer"
)

var (
	configPath *string
)

func init() {
	configPath = flag.String("config", "config.yaml", "path of config")
	flag.Parse()
}

type Span struct {
	Start   float64 `json:"start"`
	Service string  `json:"service"`
	Meta    struct {
		GRPCFullMethodName string `json:"grpc.method.name"`
	}
	Resource string `json:"resource"`
	Type     string `json:"type"`
}

func run() error {
	r, err := readJSON()
	if err != nil {
		return err
	}
	defer r.Close()

	c, err := readConfig()
	if err != nil {
		return err
	}

	spans, err := parseTrace(r)
	if err != nil {
		return err
	}

	ss := filterSpans(c, spans)
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Start < ss[j].Start
	})

	s, err := buildUML(c, ss)
	if err != nil {
		return err
	}

	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}

	return nil
}

func readConfig() (Config, error) {
	f, err := os.Open(*configPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func buildUML(config Config, spans []*Span) (string, error) {
	var b bytes.Buffer
	w := func(s string) { b.WriteString(fmt.Sprintf("%s\n", s)) }

	w("@startuml")
	for _, s := range spans {
		if s.Meta.GRPCFullMethodName == "" {
			continue
		}

		caller := s.Service
		callee := extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)
		if alias, ok := config.GRPCServiceAlias[callee]; ok {
			callee = alias
		}

		if s.Service == callee {
			continue
		}

		method := path.Base(s.Meta.GRPCFullMethodName)
		w(fmt.Sprintf(`"%s" -> "%s": %s Request`, caller, callee, method))
		w(fmt.Sprintf(`"%s" <-- "%s": %s Response`, caller, callee, method))
	}
	w("@enduml")

	return b.String(), nil
}

func extractGRPCServiceFromMethod(method string) string {
	return path.Dir(method)
}

func filterSpans(config Config, ss []*Span) []*Span {
	ss2 := ss[:0] // https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	for _, s := range ss {
		s := s
		if contains(config.IncludeServices, s.Service) &&
			!contains(config.ExcludeGRPCServices, extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)) {
			ss2 = append(ss2, s)
		}
	}
	return ss2
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func parseTrace(r io.Reader) ([]*Span, error) {
	var p interface{}
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return nil, err
	}

	rv, err := jsonpointer.Get(p, "/trace/spans")
	if err != nil {
		return nil, err
	}

	m, ok := rv.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid json")
	}

	spans := make([]*Span, 0, len(m))

	for _, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			// TODO: only skip current span
			return nil, err
		}

		var s Span
		if err := json.Unmarshal(b, &s); err != nil {
			// TODO: only skip current span
			return nil, err
		}

		spans = append(spans, &s)
	}

	return spans, nil
}

func readJSON() (io.ReadCloser, error) {
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

func main() {
	err := run()
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
