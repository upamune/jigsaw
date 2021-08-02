package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/mattn/go-jsonpointer"
)

type span struct {
	Start   float64 `json:"start"`
	Service string  `json:"service"`
	Meta    struct {
		GRPCFullMethodName string `json:"grpc.method.name"`
	}
	Resource    string   `json:"resource"`
	Type        string   `json:"type"`
	SpanID      string   `json:"span_id"`
	ParentID    string   `json:"parent_id"`
	ChildrenIDs []string `json:"children_ids"`

	children []*span
}

func resolveSpans(ss []*span) []*span {
	m := make(map[string]*span, len(ss))
	for _, s := range ss {
		s := s
		m[s.SpanID] = s
	}

	var resolve func(*span) // To call it recursively
	resolve = func(s *span) {
		if s == nil {
			return
		}
		if len(s.ChildrenIDs) == 0 {
			return
		}

		for _, cid := range s.ChildrenIDs {
			c, ok := m[cid]
			if !ok {
				continue
			}
			s.children = append(s.children, c)
			sort.Slice(s.children, func(i, j int) bool {
				return s.children[i].Start < s.children[j].Start
			})
		}

		for _, c := range s.children {
			c := c
			resolve(c)
		}
	}

	for _, s := range ss {
		if s.ParentID == "0" {
			resolve(s)
		}
	}

	spans := make([]*span, 0, len(m))
	for _, v := range m {
		if v.ParentID == "0" {
			spans = append(spans, v)
		}
	}

	sort.Slice(spans, func(i, j int) bool {
		return spans[i].Start < spans[j].Start
	})
	return spans
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
