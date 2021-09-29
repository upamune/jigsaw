package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type config struct {
	Type                string   `yaml:"type"`
	IncludeServices     []string `yaml:"include_services"`
	ExcludeGRPCServices []string `yaml:"exclude_grpc_services"`

	GRPCServiceAlias map[string]string `yaml:"grpc_serivce_alias"`
}

func readConfig(configPath string) (config, error) {
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
