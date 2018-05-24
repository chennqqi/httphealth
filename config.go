package main

type Config struct {
	ListInterval string `json:"list_interval" yaml:"list_interval"`
	ReqInterval  string `json:"req_interval" yaml:"req_interval"`
	Health       string `json:"HealthHost" yaml:"HealthHost"`
	Host         string `json:"host" yaml:"host"`
}
