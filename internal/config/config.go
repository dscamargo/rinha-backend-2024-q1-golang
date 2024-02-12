package config

import "github.com/dscamargo/rinha-2024-q1-golang/pkg"

type AppConfig struct {
	Port        string
	DatabaseUrl string
}

func Load() *AppConfig {
	return &AppConfig{
		Port:        pkg.GetEnvOrDefault("PORT", "8080"),
		DatabaseUrl: pkg.GetEnvOrDefault("DATABASE_URL", "postgresql://pg:pg@localhost:5432/rinha"),
	}
}
