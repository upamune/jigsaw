package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/upamune/jigsaw/drawer"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func createServiceNameMap(m map[string]string, spans []*span) {
	for _, s := range spans {
		s := s
		if s == nil {
			continue
		}
		m[s.SpanID] = s.Service
		if len(s.children) > 0 {
			createServiceNameMap(m, s.children)
		}
	}
}

func buildDiagram(conf config, d drawer.Drawer, spans []*span) (string, error) {
	var b bytes.Buffer
	w := func(s string) {
		if s == "" {
			return
		}
		b.WriteString(fmt.Sprintf("%s\n", s))
	}

	serviceNameMap := make(map[string]string)
	createServiceNameMap(serviceNameMap, spans)

	var draw func(s *span)
	draw = func(s *span) {
		if s == nil {
			return
		}

		caller := extractCaller(s, serviceNameMap, conf)
		callee := extractCallee(s, conf)
		isSkip := shouldSkipSpan(s, caller, callee, conf)

		method := extractMethodName(s)
		if !isSkip {
			msg := createRequestMsgWithDrawing(s, method)
			s := d.Draw(caller, callee, msg)
			w(s)
		}

		for _, c := range s.children {
			c := c
			draw(c)
		}

		if conf.NoResponse {
			return
		}
		if !isSkip {
			msg := createResponseMsgWithDrawing(s, method)
			s := d.Draw(callee, caller, msg)
			w(s)
		}
	}

	w(d.Comment("Generated by https://github.com/upamune/jigsaw"))
	w(d.Header())
	for _, s := range spans {
		draw(s)
	}
	w(d.Footer())

	return b.String(), nil
}

func convertServiceName(s string, conf config) string {
	if alias, ok := conf.ServiceAlias[s]; ok {
		return alias
	}
	return s
}

func extractCaller(s *span, serviceNameMap map[string]string, conf config) string {
	if caller, ok := serviceNameMap[s.ParentID]; ok {
		return convertServiceName(caller, conf)
	}
	return convertServiceName(s.Service, conf)
}

func extractCallee(s *span, conf config) string {
	if s.Meta.GRPCFullMethodName == "" {
		callee := extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName)
		if alias, ok := conf.GRPCServiceAlias[callee]; ok {
			callee = alias
		}
	}
	return convertServiceName(s.Service, conf)
}

func extractMethodName(s *span) string {
	if s.Meta.GRPCFullMethodName != "" {
		return path.Base(s.Meta.GRPCFullMethodName)
	}
	if s.Type == "graphql" {
		name, err := extractGraphQLOperationName(s.Resource)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse graphql query(%s): %v\n", s.Resource, err)
			return "GraphQL"
		}
		return fmt.Sprintf("GraphQL: %s", name)
	}
	return s.Resource
}

func shouldSkipSpan(s *span, caller, callee string, conf config) bool {
	if conf.IsSkipSelfCall {
		return caller == callee
	}
	if len(conf.IncludeServices) > 0 {
		fmt.Println(conf.IncludeServices)
		return !contains(conf.IncludeServices, s.Service)
	}
	return contains(conf.ExcludeGRPCServices, extractGRPCServiceFromMethod(s.Meta.GRPCFullMethodName))
}

func extractGRPCServiceFromMethod(method string) string {
	return path.Dir(method)
}

func extractGraphQLOperationName(query string) (string, error) {
	doc, err := parser.ParseQuery(&ast.Source{Input: query})
	if err != nil {
		return "", err
	}
	var names []string
	for _, op := range doc.Operations {
		names = append(names, op.Name)
	}
	return strings.Join(names, ","), nil
}

func createRequestMsgWithDrawing(s *span, method string) string {
	if s.Type == "" || s.Type == "sql" {
		return method
	}
	return fmt.Sprintf("%s Request", method)
}

func createResponseMsgWithDrawing(s *span, method string) string {
	if s.Type == "" {
		return method
	}
	if s.Type == "sql" {
		return ""
	}
	return fmt.Sprintf("%s Response", method)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
