package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/subosito/gotenv"
)

type Config struct {
	App    AppConfig
	Server ServerConfig
	DB     DBConfig
}

type AppConfig struct {
	Env     string
	BaseURL string
}

type ServerConfig struct {
	Port         string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	IdleTime     time.Duration
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (db *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

func Load() (*Config, error) {
	_ = gotenv.Load(".env")

	cfg := &Config{
		App: AppConfig{
			Env:     getStr("APP_ENV", "development"),
			BaseURL: getStr("BASE_URL", "http://localhost:8080"),
		},

		Server: ServerConfig{
			Port:         getStr("PORT", "8080"),
			ReadTimeOut:  getDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeOut: getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTime:     getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},

		DB: DBConfig{
			Host:            getStr("DB_HOST", "localhost"),
			Port:            getStr("DB_PORT", "5432"),
			User:            getStr("DB_USER", "postgres"),
			Password:        getStr("DB_PASSWORD", ""),
			Name:            getStr("DB_NAME", ""),
			SSLMode:         getStr("DB_SSLMODE", "disable"),
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
	}

	return cfg, nil
}

func getStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
