package config

import (
	"fmt"
	"os"
	"strconv"
)

const defaultPort = 8080

type ServiceConfig struct {
	Port int
}

func LoadServiceConfig() (ServiceConfig, error) {
	port := defaultPort
	if raw := os.Getenv("PORT"); raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil {
			return ServiceConfig{}, fmt.Errorf("invalid PORT %q: %w", raw, err)
		}
		port = p
	}
	return ServiceConfig{Port: port}, nil
}

func (c ServiceConfig) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
