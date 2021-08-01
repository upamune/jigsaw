package main

type Config struct {
	IncludeServices     []string `yaml:"include_services"`
	ExcludeGRPCServices []string `yaml:"exclude_grpc_services"`

	GRPCServiceAlias map[string]string `yaml:"grpc_serivce_alias"`
}
