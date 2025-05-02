package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Telegram TelegramConfig
	Postgres PostgresConfig
}

type PostgresConfig struct {
	User   string
	Pass   string
	Dbname string
	Port   string
	Host   string
}

type TelegramConfig struct {
	Token string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()

	return &Config{
		Telegram: TelegramConfig{
			Token: viper.GetString("TELEGRAM_BOT_TOKEN"),
		},
		Postgres: PostgresConfig{
			User:   viper.GetString("DB_USERNAME"),
			Pass:   viper.GetString("DB_PASSWORD"),
			Dbname: viper.GetString("DB_DATABASE"),
			Port:   viper.GetString("DB_PORT"),
			Host:   viper.GetString("DB_HOST"),
		},
	}, nil
}
