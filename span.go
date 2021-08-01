package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/mattn/go-jsonpointer"
)

type span struct {
	Start   float64 `json:"start"`
	Service string  `json:"service"`
	Meta    struct {
		GRPCFullMethodName string `json:"grpc.method.name"`
	}
	Resource string `json:"resource"`
	Type     string `json:"type"`
}

func filterSpans(conf config, ss []*span) []*span {
	ss2 := ss[:0] // https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	for _, s := range ss {
		s := s
		if contains(conf.IncludeServices, s.Service) &&
			!contains(conf.ExcludeGRPCServices, extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)) {
			ss2 = append(ss2, s)
		}
	}
	return ss2
}

func parseSpans(r io.Reader) ([]*span, error) {
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

	spans := make([]*span, 0, len(m))

	for _, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			// TODO: only skip current span
			return nil, err
		}

		var s span
		if err := json.Unmarshal(b, &s); err != nil {
			// TODO: only skip current span
			return nil, err
		}

		spans = append(spans, &s)
	}

	return spans, nil
}
