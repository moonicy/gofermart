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

func (c *Config) parseFlag() {
	flag.StringVar(&c.Host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.DatabaseURI, "d", "host=localhost port=5432 user=mila password=qwerty dbname=gofermart sslmode=disable", "database uri")
	flag.StringVar(&c.AccrualSystemAddress, "r", "http://localhost:8081", "accrual system address")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		c.Host = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		c.DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		c.AccrualSystemAddress = envAccrualSystemAddress
	}
}
