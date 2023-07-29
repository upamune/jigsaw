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

	GRPCServiceAlias map[string]string `yaml:"grpc_service_alias"`
	// GRPCSerivceAlias has a typo, but it's kept for the sake of backward compatibility.
	GRPCSerivceAlias map[string]string `yaml:"grpc_serivce_alias"`
	ServiceAlias     map[string]string `yaml:"service_alias"`

	IsSkipSelfCall bool `yaml:"is_skip_self_call"`
	NoResponse     bool `yaml:"no_response"`

	Debug bool `yaml:"debug"`
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

	// NOTE: The GRPCSerivceAlias contains a typo, but in order to correctly load settings that have past typos,
	// we only overwrite it when GRPCServiceAlias is empty.
	if len(c.GRPCServiceAlias) == 0 && len(c.GRPCSerivceAlias) > 0 {
		c.GRPCServiceAlias = c.GRPCSerivceAlias
	}

	return c, nil
}
