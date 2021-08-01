package main

import (
	"bytes"
	"fmt"
	"path"
)

func buildUML(conf config, spans []*span) (string, error) {
	var b bytes.Buffer
	w := func(s string) { b.WriteString(fmt.Sprintf("%s\n", s)) }

	w("@startuml")
	for _, s := range spans {
		if s.Meta.GRPCFullMethodName == "" {
			continue
		}

		caller := s.Service
		callee := extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)
		if alias, ok := conf.GRPCServiceAlias[callee]; ok {
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
