package config

import "github.com/coreyvan/kid-dictionary/pkg/env"

type Config struct {
	BindAddr  string
	Port      string
	PrettyLog bool
	LogLevel  string

	// OpenAI configuration
	OpenAIAPIKey string
	OpenAIModel  string
}

func GetConfig() Config {
	return Config{
		BindAddr:  env.GetString("BIND_ADDR", "0.0.0.0"),
		Port:      env.GetString("LISTEN_PORT", "8080"),
		PrettyLog: env.GetBool("PRETTY_LOG", false),
		LogLevel:  env.GetString("LOG_LEVEL", "info"),

		OpenAIAPIKey: env.GetString("OPENAI_API_KEY", ""),
		OpenAIModel:  env.GetString("OPENAI_MODEL", "gpt-4o-mini"),
	}
}
