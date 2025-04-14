package config

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// if envVar, ok := os.LookupEnv("RUN_ADDRESS"); ok && envVar != "" {
// 	config.RunAddress = envVar
// }
// if envVar, ok := os.LookupEnv("DATABASE_URI"); ok && envVar != "" {
// 	config.DatabaseURI = envVar
// }
// if envVar, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok && envVar != "" {
// 	config.AccrualSystemAddress = envVar
// }
// if envVar, ok := os.LookupEnv("SECRET_KEY"); ok && envVar != "" {
// 	config.SecretKey = envVar
// }

// if envVar, ok := os.LookupEnv("TOKEN_VALIDITY"); ok && envVar != "" {

// 	duration, err := time.ParseDuration(envVar)
// 	if err != nil {
// 		panic(err)
// 	}
// 	config.TokenValidityDuration = duration
// }

func TestParseEnv(t *testing.T) {

	// Test cases
	tests := []struct {
		name                 string
		runAddress           string
		databaseURI          string
		accrualSystemAddress string
		secretKey            string
		tokenValidity        string
		expected             *Config
	}{
		{"Test1", ":8080", "uri", ":9001", "secretkey", "1m", &Config{":8080", "uri", ":9001", "secretkey", 1 * time.Minute}},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			oldRunAddress := os.Getenv("RUN_ADDRESS")
			oldDatabaseURI := os.Getenv("DATABASE_URI")
			oldAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
			oldSecretKey := os.Getenv("SECRET_KEY")
			oldTokenValidity := os.Getenv("TOKEN_VALIDITY")

			if err := os.Setenv("RUN_ADDRESS", tt.runAddress); err != nil {
				panic(err)
			}

			if err := os.Setenv("DATABASE_URI", tt.databaseURI); err != nil {
				panic(err)
			}

			if err := os.Setenv("ACCRUAL_SYSTEM_ADDRESS", tt.accrualSystemAddress); err != nil {
				panic(err)
			}

			if err := os.Setenv("SECRET_KEY", tt.secretKey); err != nil {
				panic(err)
			}

			if err := os.Setenv("TOKEN_VALIDITY", tt.tokenValidity); err != nil {
				panic(err)
			}

			config := &Config{}
			parseEnv(config)

			if err := os.Setenv("RUN_ADDRESS", oldRunAddress); err != nil {
				panic(err)
			}
			if err := os.Setenv("DATABASE_URI", oldDatabaseURI); err != nil {
				panic(err)
			}
			if err := os.Setenv("ACCRUAL_SYSTEM_ADDRESS", oldAccrualSystemAddress); err != nil {
				panic(err)
			}
			if err := os.Setenv("SECRET_KEY", oldSecretKey); err != nil {
				panic(err)
			}
			if err := os.Setenv("TOKEN_VALIDITY", oldTokenValidity); err != nil {
				panic(err)
			}

			if diff := cmp.Diff(config, tt.expected); diff != "" {
				t.Errorf("Structs mismatch (-config +expected):\n%s", diff)
			}
		})
	}
}
