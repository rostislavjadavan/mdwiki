package config

import "github.com/ilyakaznacheev/cleanenv"

type AppConfig struct {
	Port    string `yaml:"port" env:"PORT" env-default:"8080"`
	Host    string `yaml:"host" env:"HOST" env-default:"localhost"`
	Storage string `yaml:"storage" env:"STORAGE" env-default:".storage"`
	Theme   string `yaml:"theme" env:"THEME" env-default:"system"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var config AppConfig
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
