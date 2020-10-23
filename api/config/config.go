package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

// Conf is a the main config struct.
type Conf struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConf
	Db     dbConf
}

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

type serverConf struct {
	Port         int           `env:"PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

// AppConfig is a function to create a config struct and read env variables into it.
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
