package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Host         string
	HttpPort     int
	NodeID       uint16
	Logging      Logging
	Database     Database
	Cache        Cache
	MessageQueue MessageQueue
	Token        Token
	Session      Session
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
	Key    string
}

type Session struct {
	Exp        time.Duration
	RefreshExp time.Duration
}

func LoadConfig() *AppConfig {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	return &AppConfig{
		Host:     viper.GetString("HOSTNAME"),
		HttpPort: viper.GetInt("HTTP_PORT"),
		NodeID:   viper.GetUint16("NODE_ID"),
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
			Key:    viper.GetString("JWT_KEY"),
		},
		Session: Session{
			Exp:        viper.GetDuration("SESSION_EXP"),
			RefreshExp: viper.GetDuration("SESSION_REFRESH_EXP"),
		},
	}
}
