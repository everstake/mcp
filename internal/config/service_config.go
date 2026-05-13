package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

const (
	TransportHTTP  = "http"
	TransportStdio = "stdio"
)

type ServiceConfig struct {
	DashboardURL    string `env:"DASHBOARD_URL,required" envDefault:"https://dashboard-api.everstake.one"`
	DashboardAPIKey string `env:"DASHBOARD_API_KEY"`
	Transport       string `env:"MCP_TRANSPORT" envDefault:"http"`
	Port            int    `env:"PORT" envDefault:"8080"`
}

func LoadServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	if err := env.Parse(&cfg); err != nil {
		return ServiceConfig{}, fmt.Errorf("load config from env: %w", err)
	}
	switch cfg.Transport {
	case TransportHTTP, TransportStdio:
	default:
		return ServiceConfig{}, fmt.Errorf("invalid MCP_TRANSPORT %q: must be %q or %q", cfg.Transport, TransportHTTP, TransportStdio)
	}
	return cfg, nil
}
