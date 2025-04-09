package config

import (
	"os"
)

func parseEnv(config *Config) {
	if envVar, ok := os.LookupEnv("RUN_ADDRESS"); ok && envVar != "" {
		config.RunAddress = envVar
	}
	if envVar, ok := os.LookupEnv("DATABASE_URI"); ok && envVar != "" {
		config.DatabaseURI = envVar
	}
	if envVar, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok && envVar != "" {
		config.AccrualSystemAddress = envVar
	}
}
