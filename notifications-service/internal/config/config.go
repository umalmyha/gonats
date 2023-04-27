package config

import "github.com/caarlos0/env/v8"

type NATS struct {
	URL        string `env:"NATS_URL,notEmpty"`
	StreamName string `env:"NATS_STREAM_NAME,notEmpty"`
	Username   string `env:"NATS_USERNAME"`
	Password   string `env:"NATS_PASSWORD"`
}

type Config struct {
	NATS
}

func Build() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
