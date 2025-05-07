package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Kafka    Kafka
	Http     Http
	Services Services
}

type Kafka struct {
	Topic   string   `env:"KAFKA_TOPIC"`
	GroupId string   `env:"KAFKA_GROUP_ID"`
	Brokers []string `env:"KAFKA_BROKERS"`
}

type Http struct {
	Host        string `env:"HTTP_HOST"`
	Port        uint16 `env:"HTTP_PORT"`
	MetricsPort uint16 `env:"HTTP_METRICS_PORT"`
}

type Services struct {
	ApplicationLevelBaseUrl string `env:"APPLICATION_LEVEL_BASE_URL"`
}

func Load() (Config, error) {
	var config Config
	err := env.Parse(&config)
	return config, err
}
