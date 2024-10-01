package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (c *DBConfig) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *DBConfig) Migrate() {
	m, err := migrate.New(
		"file://internal/database/migrations",
		c.DSN(),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		panic(err)
	}
}
