package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Mongo Mongo `yaml:"mongo"`
	}

	Mongo struct {
		Url      string `env:"MONGO_URL" yaml:"url"`
		Database string `env:"MONGO_DATABASE" yaml:"database"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("config/config.yml", cfg); err != nil {
		return nil, err
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
