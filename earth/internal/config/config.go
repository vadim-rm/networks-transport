package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Http     Http
	Services Services
}

type Http struct {
	Host        string `env:"HTTP_HOST"`
	Port        uint16 `env:"HTTP_PORT"`
	MetricsPort uint16 `env:"HTTP_METRICS_PORT"`
}

type Services struct {
	DataLinkLevelBaseUrl string `env:"DATA_LINK_LEVEL_BASE_URL"`
}

func Load() (Config, error) {
	var config Config
	err := env.Parse(&config)
	return config, err
}
