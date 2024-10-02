package utils

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     int    `mapstructure:"DB_PORT"`
	SSLMode    string `mapstructure:"DB_SSL_MODE"`
	TimeZone   string `mapstructure:"DB_TIMEZONE"`
}

func NewConfig(env string) *Config {
	config := &Config{}

	config.loadConfig(env)

	return config
}

func (c *Config) loadConfig(env string) (err error) {
	path, _ := os.Getwd()

	viper.AddConfigPath(path)
	viper.SetConfigName("app-" + env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
