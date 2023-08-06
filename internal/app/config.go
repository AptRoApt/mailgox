package app

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	ImapPort int    `yaml:"imap-port"`
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	return cfg, err
}
