package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type ServiceConfig struct {
	DashboardURL    string `env:"DASHBOARD_URL,required" envDefault:"https://dashboard-api.everstake.one"`
	DashboardAPIKey string `env:"DASHBOARD_API_KEY"`
	Port            int    `env:"PORT" envDefault:"8080"`
}

func LoadServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	if err := env.Parse(&cfg); err != nil {
		return ServiceConfig{}, fmt.Errorf("load config from env: %w", err)
	}
	return cfg, nil
}
