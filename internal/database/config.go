package database

import (
	"context"
	"fmt"
	"time"

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
	minConns int32
	maxConns int32
}

func NewDBConfig(
	host string,
	user string,
	password string,
	dbName string,
	port int,
	sslMode string,
	timeZone string,
	minConns int32,
	maxConns int32,
) *DBConfig {
	return &DBConfig{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbName,
		Port:     port,
		SSLMode:  sslMode,
		TimeZone: timeZone,
		minConns: minConns,
		maxConns: maxConns,
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode, c.TimeZone)
}

func (c *DBConfig) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(c.DSN())
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = c.maxConns
	dbConfig.MinConns = c.minConns
	dbConfig.MaxConnIdleTime = time.Minute * 10

	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
