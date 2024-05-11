package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultCount = 1

type LoggerConfig struct {
	LogFile  string
	MaxCount int
}

type DBConfig struct {
	PgUser    string
	PgPass    string
	PgHost    string
	PgPort    uint16
	PgDb      string
	PgSSLMode string
}

type ServerConfig struct {
	HTTPPort       string
	ServerEndpoint string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func GetLoggerConfig() LoggerConfig {
	maxCount, err := strconv.Atoi(os.Getenv("MAXCOUNT"))
	if err != nil {
		// log: Failed to interpritate number from .env so use 1
		maxCount = defaultCount
	}
	return LoggerConfig{
		LogFile:  os.Getenv("LOGFILE"),
		MaxCount: maxCount,
	}
}

func GetDBConfig() (DBConfig, error) {
	pgPort, err := strconv.ParseInt(os.Getenv("PGPORT"), 0, 16)
	if err != nil {
		return DBConfig{}, err
	}

	return DBConfig{
		PgUser:    os.Getenv("PGUSER"),
		PgPass:    os.Getenv("PGPASSWORD"),
		PgHost:    os.Getenv("PGHOST"),
		PgPort:    uint16(pgPort),
		PgDb:      os.Getenv("PGDATABASE"),
		PgSSLMode: os.Getenv("PGSSLMODE"),
	}, nil
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort:       ":" + os.Getenv("HTTP_PORT"),
		ServerEndpoint: os.Getenv("SERVER_ENDPOINT"),
	}
}
