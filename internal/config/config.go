package config

import "time"

type Config struct {
	RunAddress            string
	DatabaseURI           string
	AccrualSystemAddress  string
	SecretKey             string
	TokenValidityDuration time.Duration
}

func ParseConfig() (*Config, error) {
	config := &Config{}
	parseFlags(config)
	parseEnv(config)
	return config, nil
}
