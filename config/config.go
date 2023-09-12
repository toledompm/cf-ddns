package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Cloudflare struct {
		// Token is the API token for Cloudflare
		Token string `json:"token"`
		// ZoneID is the zone ID where all records are located
		ZoneID string `json:"zoneId"`
	} `json:"cloudflare"`

	Records []struct {
		// Name is the name of the record
		Name string `json:"name"`
		// Proxy is a flag to enable or disable proxy
		Proxy bool `json:"proxy"`
	} `json:"records"`

	IPV6 struct {
		// Enabled is a flag to enable or disable IPV6
		Enabled bool `json:"enabled"`
		// FetchAddress is the URL to fetch the current IPV6 address
		FetchAddress string `json:"fetchAddress"`
	} `json:"ipv6"`

	IPV4 struct {
		// Enabled is a flag to enable or disable IPV4
		Enabled bool `json:"enabled"`
		// FetchAddress is the URL to fetch the current IPV4 address
		FetchAddress string `json:"fetchAddress"`
	} `json:"ipv4"`
}

func New() *Config {
	return &Config{}
}

func MustParseConfig(jsonBytes []byte, cfg *Config) *Config {
	err := json.Unmarshal(jsonBytes, cfg)

	if err != nil {
		fmt.Println("Error parsing config file")
		panic(err)
	}

	return cfg
}
