package config

type Config struct {
	SecretKey string
}

func NewConfig() *Config {
	return &Config{SecretKey: "very secret"}
}
