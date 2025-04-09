package config

import (
	"flag"
	"time"
)

func parseFlags(config *Config) {

	flag.StringVar(&config.RunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.SecretKey, "k", "secretKey", "jwt token signing key")
	flag.DurationVar(&config.TokenValidityDuration, "v", 5*time.Minute, "jwt token validity duration time interval")
	flag.StringVar(&config.DatabaseURI, "d", "", "database URI")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "accrual system address")
	flag.Parse()

}
