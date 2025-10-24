package app

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	Host      string
	JWTSecret string
	LogLevel  string
	Database  DatabaseConfig
	HTTPS     HTTPSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type HTTPSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

func LoadConfig() *Config {
	config := &Config{
		Port:      getEnv("PORT", "8080"),
		Host:      getEnv("HOST", "0.0.0.0"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "vkarmane"),
			Password: getEnv("DB_PASSWORD", "vkarmane_password"),
			DBName:   getEnv("DB_NAME", "vkarmane"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		HTTPS: HTTPSConfig{
			Enabled:  getEnv("HTTPS_ENABLED", "true") == "true",
			CertFile: getEnv("HTTPS_CERT_FILE", "ssl/server.crt"),
			KeyFile:  getEnv("HTTPS_KEY_FILE", "ssl/server.key"),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetCORSOrigins() []string {
	origins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:8080,https://localhost:3000,https://localhost:8080")
	if origins == "" {
		return []string{}
	}

	var result []string
	for _, origin := range []string{origins} {
		if origin != "" {
			result = append(result, origin)
		}
	}
	return result
}

func (c *Config) IsProduction() bool {
	return getEnv("ENV", "development") == "production"
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
