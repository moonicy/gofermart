package config

import (
	"flag"
	"os"
)

type Config struct {
	Host                 string
	DatabaseURI          string
	AccrualSystemAddress string
}

func NewConfig() Config {
	sc := Config{}
	sc.parseFlag()
	return sc
}

func (sc *Config) parseFlag() {
	flag.StringVar(&sc.Host, "a", "localhost:8081", "address and port to run server")
	flag.StringVar(&sc.DatabaseURI, "d", "host=localhost port=5432 user=mila password=qwerty dbname=gofermart sslmode=disable", "database uri")
	flag.StringVar(&sc.AccrualSystemAddress, "r", "http://localhost:8080", "accrual system address")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		sc.Host = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		sc.DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		sc.AccrualSystemAddress = envAccrualSystemAddress
	}
}
