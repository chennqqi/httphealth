package main

type Config struct {
	ListInterval string `json:"list_interval" yaml:"list_interval"`
	ReqInterval  string `json:"req_interval" yaml:"req_interval"`
	HealthHost   string `json:"health" yaml:"health"`
}
