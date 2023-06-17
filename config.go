package main

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/upamune/jigsaw/drawer"
)

type config struct {
	Type                string   `yaml:"type"`
	IncludeServices     []string `yaml:"include_services"`
	ExcludeGRPCServices []string `yaml:"exclude_grpc_services"`

	GRPCServiceAlias map[string]string `yaml:"grpc_serivce_alias"`
	ServiceAlias     map[string]string `yaml:"service_alias"`

	IsSkipSelfCall bool `yaml:"is_skip_self_call"`

	NoResponse bool `yaml:"no_response"`
}

var defaultConfig = config{
	Type:           drawer.TypeMermaid,
	IsSkipSelfCall: true,
	NoResponse:     false,
}

func readConfig(configPath string) (config, error) {
	if configPath == "" {
		return defaultConfig, nil
	}
	f, err := os.Open(configPath)
	if err != nil {
		return config{}, err
	}
	defer f.Close()

	var c config
	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return config{}, err
	}

	return c, nil
}
