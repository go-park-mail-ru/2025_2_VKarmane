package app

import (
	"os"
)

type Config struct {
	Port      string
	Host      string
	JWTSecret string
	LogLevel  string
}

func LoadConfig() *Config {
	config := &Config{
		Port:      getEnv("PORT", "8080"),
		Host:      getEnv("HOST", "0.0.0.0"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

func (c *Config) GetServerAddress() string {
	return c.Host + ":" + c.Port
}

func (c *Config) GetCORSOrigins() []string {
	if c.IsProduction() {
		return []string{
			"http://217.16.23.67:8000",
			"https://217.16.23.67:8000",
		}
	}
	return []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
		"http://217.16.23.67:8000",
	}
}
