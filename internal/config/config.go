package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type LoggerConfig struct {
	Enviroment string
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
	Timeout        time.Duration
	IdleTimeout    time.Duration
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
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

func GetServerConfig() (ServerConfig, error) {
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		return ServerConfig{}, err
	}

	idletimeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		HTTPPort:       ":" + os.Getenv("HTTP_PORT"),
		ServerEndpoint: os.Getenv("SERVER_ENDPOINT"),
		Timeout:        timeout,
		IdleTimeout:    idletimeout,
	}, nil
}
