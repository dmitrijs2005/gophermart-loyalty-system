package config

import (
	"os"
	"time"
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
	if envVar, ok := os.LookupEnv("SECRET_KEY"); ok && envVar != "" {
		config.SecretKey = envVar
	}

	if envVar, ok := os.LookupEnv("TOKEN_VALIDITY"); ok && envVar != "" {

		duration, err := time.ParseDuration(envVar)
		if err != nil {
			panic(err)
		}
		config.TokenValidityDuration = duration
	}

}
