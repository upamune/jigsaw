package main

import (
	"bytes"
	"fmt"
	"path"
)

func buildUML(conf config, spans []*span) (string, error) {
	var b bytes.Buffer
	w := func(s string) { b.WriteString(fmt.Sprintf("%s\n", s)) }

	var draw func(s *span)
	draw = func(s *span) {
		if s == nil {
			return
		}

		caller := s.Service
		callee := extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)
		if alias, ok := conf.GRPCServiceAlias[callee]; ok {
			callee = alias
		}

		isSkip := s.Meta.GRPCFullMethodName == "" ||
			s.Service == callee ||
			!contains(conf.IncludeServices, s.Service) ||
			contains(conf.ExcludeGRPCServices, extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName))

		method := path.Base(s.Meta.GRPCFullMethodName)
		if !isSkip {
			w(fmt.Sprintf(`"%s" -> "%s": %s Request`, caller, callee, method))
		}

		for _, c := range s.children {
			c := c
			draw(c)
		}

		if *noResponse {
			return
		}
		if !isSkip {
			w(fmt.Sprintf(`"%s" <-- "%s": %s Response`, caller, callee, method))
		}
	}

	w("@startuml")
	for _, s := range spans {
		draw(s)
	}
	w("@enduml")

	return b.String(), nil
}
