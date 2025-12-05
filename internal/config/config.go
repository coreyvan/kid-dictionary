package config

import (
	"fmt"

	"github.com/coreyvan/kid-dictionary/pkg/env"
)

type Config struct {
	BindAddr  string
	Port      string
	PrettyLog bool
	LogLevel  string

	// OpenAI configuration
	OpenAIAPIKey string
	OpenAIModel  string

	// Database configuration
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
}

func GetConfig() Config {
	cfg := Config{
		BindAddr:  env.GetString("BIND_ADDR", "0.0.0.0"),
		Port:      env.GetString("LISTEN_PORT", "8080"),
		PrettyLog: env.GetBool("PRETTY_LOG", false),
		LogLevel:  env.GetString("LOG_LEVEL", "info"),

		OpenAIAPIKey: env.GetString("OPENAI_API_KEY", ""),
		OpenAIModel:  env.GetString("OPENAI_MODEL", "gpt-4o-mini"),

		DatabaseURL: env.GetString("DATABASE_URL", ""),
		DBHost:      env.GetString("DB_HOST", "localhost"),
		DBPort:      env.GetString("DB_PORT", "5432"),
		DBUser:      env.GetString("DB_USER", "postgres"),
		DBPassword:  env.GetString("DB_PASSWORD", ""),
		DBName:      env.GetString("DB_NAME", "kid_dictionary"),
	}

	// Build DATABASE_URL from components if not explicitly set
	if cfg.DatabaseURL == "" && cfg.DBHost != "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	}

	return cfg
}
