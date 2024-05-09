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

func GetLoggerConfig() LoggerConfig {
	maxCount, err := strconv.Atoi(getVarFromEnv("MAXCOUNT"))
	if err != nil {
		// log: Failed to interpritate number from .env so use 1
		maxCount = defaultCount
	}
	return LoggerConfig{
		LogFile:  getVarFromEnv("LOGFILE"),
		MaxCount: maxCount,
	}
}

func GetDBConfig() (DBConfig, error) {
	pgPort, err := strconv.ParseInt(getVarFromEnv("PGPORT"), 0, 16)
	if err != nil {
		return DBConfig{}, err
	}

	return DBConfig{
		PgUser:    getVarFromEnv("PGUSER"),
		PgPass:    getVarFromEnv("PGPASSWORD"),
		PgHost:    getVarFromEnv("PGHOST"),
		PgPort:    uint16(pgPort),
		PgDb:      getVarFromEnv("PGDATABASE"),
		PgSSLMode: getVarFromEnv("PGSSLMODE"),
	}, nil
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort:       ":" + getVarFromEnv("HTTP_PORT"),
		ServerEndpoint: getVarFromEnv("SERVER_ENDPOINT"),
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func getVarFromEnv(varName string) string {
	return os.Getenv(varName)
}
