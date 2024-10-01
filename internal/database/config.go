package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SSLMode  string
	TimeZone string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		DBName:   "url_shortener",
		Port:     5433,
		SSLMode:  "disable",
		TimeZone: "America/Sao_Paulo",
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode, c.TimeZone)
}

func (c *DBConfig) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, c.DSN())
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
