package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	HttpPort     int
	Logging      Logging
	Database     Database
	Cache        Cache
	MessageQueue MessageQueue
	Token        Token
}

type Logging struct {
	Level string
}

type Database struct {
	Address  string
	Username string
	Password string
	Database string
}

func (d *Database) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.Username, d.Password, d.Address, d.Database)
}

type Cache struct {
	Address string
}

type MessageQueue struct {
	Brokers string
}

func (m *MessageQueue) GetBrokerList() []string {
	return strings.Split(m.Brokers, ",")
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
		HttpPort: viper.GetInt("HTTP_PORT"),
		Logging: Logging{
			Level: viper.GetString("LOG_LEVEL"),
		},
		Database: Database{
			Address:  viper.GetString("DATABASE_ADDRESS"),
			Username: viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Database: viper.GetString("DATABASE_DB"),
		},
		Cache: Cache{
			Address: viper.GetString("CACHE_ADDRESS"),
		},
		MessageQueue: MessageQueue{
			Brokers: viper.GetString("KAFKA_BROKERS"),
		},
		Token: Token{
			Secret: viper.GetString("JWT_SECRET"),
			Exp:    viper.GetDuration("JWT_EXP"),
		},
	}
}
