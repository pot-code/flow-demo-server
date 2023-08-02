package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port     int
	Logging  *Logging
	Database *Database
	Cache    *Cache
	Token    *Token
}

type Logging struct {
	Level string
}

type Database struct {
	DSN string
}

type Cache struct {
	DSN string
}

type Token struct {
	Secret string
	Exp    time.Duration
}

func LoadConfig() *AppConfig {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	return &AppConfig{
		Port: viper.GetInt("PORT"),
		Logging: &Logging{
			Level: viper.GetString("LOG_LEVEL"),
		},
		Database: &Database{
			DSN: viper.GetString("DATABASE_DSN"),
		},
		Cache: &Cache{
			DSN: viper.GetString("CACHE_DSN"),
		},
		Token: &Token{
			Secret: viper.GetString("JWT_SECRET"),
			Exp:    viper.GetDuration("JWT_EXP"),
		},
	}
}
