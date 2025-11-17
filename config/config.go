package config

import "github.com/caarlos0/env/v11"

type (
	Config struct {
		HTTP HTTP
		PG   PG
	}

	HTTP struct {
		URL  string `env:"HTTP_URL" envDefault:"localhost:8080"`
		PORT string `env:"HTTP_PORT" envDefault:"8080"`
	}

	PG struct {
		URL string `env:"PG_URL,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
