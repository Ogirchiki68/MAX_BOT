package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN,required"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
