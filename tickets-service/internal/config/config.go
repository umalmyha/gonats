package config

import "github.com/caarlos0/env/v8"

type NATS struct {
	URL                    string `env:"NATS_URL,notEmpty"`
	NotificationStreamName string `env:"NATS_NOTIFICATION_STREAM_NAME,notEmpty"`
	Username               string `env:"NATS_USERNAME"`
	Password               string `env:"NATS_PASSWORD"`
}

type Config struct {
	Port        int    `env:"PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,notEmpty"`
	NATS        NATS
}

func Build() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
