package utils

import (
	"log"
	"os"
	"path/filepath"

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
	path, _ := GetRootPath()

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

// GetRootPath returns the root path of the project
func GetRootPath() (ex string, err error) {
	ex, _ = os.Getwd()
	fileToStat := "go.mod"
	if os.Getenv("ENV") == "prod" {
		fileToStat = "main"
	}

	_, err = os.Stat(filepath.Join(ex, fileToStat))

	if err != nil {
		for i := 0; i < 5; i++ {
			ex = filepath.Join(ex, "../")
			_, err = os.Stat(filepath.Join(ex, fileToStat))

			if err == nil {
				break
			}
		}
		if err != nil {
			log.Println("No env file provided, using only env variables")
		}
	}
	return
}
