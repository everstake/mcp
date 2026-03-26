package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type ServiceConfig struct {
	Port         int    `env:"PORT" envDefault:"8080"`
	DashboardUrl string `env:"DASHBOARD_URL,required" envDefault:"https://dashboard-api.everstake.one"`
}

func LoadServiceConfig() (ServiceConfig, error) {
	cfg := ServiceConfig{}
	if err := env.Parse(&cfg); err != nil {
		return ServiceConfig{}, fmt.Errorf("load config from env: %w", err)
	}
	return cfg, nil
}
