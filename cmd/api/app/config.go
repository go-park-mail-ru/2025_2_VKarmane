package app

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port               string
	AuthServicePort    string
	AuthServiceHost    string
	BudgetServicePort  string
	BudgetServiceHost  string
	FinanceServicePort string
	FinanceServiceHost string
	KafkaProducerHost  string
	KafkaProducerPort  string
	Host               string
	JWTSecret          string
	LogLevel           string
	Database           DatabaseConfig
	HTTPS              HTTPSConfig
	MinIO              MinIOConfig
	ElasticSearch      ElasticSearchConfig
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

type MinIOConfig struct {
	Endpoint   string
	Port       string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

type ElasticSearchConfig struct {
	Host string
	Port string
}

func LoadConfig() *Config {
	config := &Config{
		Port:               getEnv("PORT", "8080"),
		AuthServicePort:    getEnv("AUTH_SERVICE_PORT", "8090"),
		AuthServiceHost:    getEnv("AUTH_SERVICE_HOST", "auth_service"),
		BudgetServicePort:  getEnv("BUDGET_SERVICE_PORT", "8100"),
		BudgetServiceHost:  getEnv("BUDGET_SERVICE_HOST", "budget_service"),
		FinanceServicePort: getEnv("FINANCE_SERVICE_PORT", "8110"),
		FinanceServiceHost: getEnv("FINANCE_SERVICE_HOST", "finance_service"),
		KafkaProducerHost:  getEnv("KAFKA_PRODUCER_HOST", "kafka"),
		KafkaProducerPort:  getEnv("KAFKA_PRODUCER_PORT", "9092"),
		Host:               getEnv("HOST", "0.0.0.0"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "vkarmane"),
			Password: getEnv("DB_PASSWORD", "vkarmane_password"),
			DBName:   getEnv("DB_NAME", "vkarmane"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		HTTPS: HTTPSConfig{
			Enabled:  getEnv("HTTPS_ENABLED", "false") == "true",
			CertFile: getEnv("HTTPS_CERT_FILE", "ssl/server.crt"),
			KeyFile:  getEnv("HTTPS_KEY_FILE", "ssl/server.key"),
		},
		MinIO: MinIOConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost"),
			Port:       getEnv("MINIO_PORT", "9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin123"),
			UseSSL:     getEnv("MINIO_USE_SSL", "true") == "true",
			BucketName: getEnv("MINIO_BUCKET_NAME", "images"),
		},
		ElasticSearch: ElasticSearchConfig{
			Port: getEnv("ELASTIC_SEARCH_PORT", "9200"),
			Host: getEnv("ELASTIC_SEARCH_HOST", "elasticsearch"),
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
	var result []string

	corsHost := getEnv("CORS_HOST", "localhost")
	corsFrontendPort := getEnv("CORS_FRONTEND_PORT", "8000")

	if corsHost != "" {
		if corsFrontendPort != "" {
			result = append(result, fmt.Sprintf("http://%s:%s", corsHost, corsFrontendPort))
			result = append(result, fmt.Sprintf("https://%s:%s", corsHost, corsFrontendPort))
		} else {
			result = append(result, fmt.Sprintf("http://%s", corsHost))
			result = append(result, fmt.Sprintf("https://%s", corsHost))
		}
	}

	serverHost := c.Host
	if serverHost == "0.0.0.0" {
		serverHost = "localhost"
	}
	serverPort := c.Port

	result = append(result, fmt.Sprintf("http://%s:%s", serverHost, serverPort))
	if c.HTTPS.Enabled {
		result = append(result, fmt.Sprintf("https://%s:%s", serverHost, serverPort))
	}

	corsOriginsEnv := getEnv("CORS_ORIGINS", "")
	if corsOriginsEnv != "" {
		origins := strings.Split(corsOriginsEnv, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				result = append(result, origin)
			}
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

func (c *Config) GetMinIOEndpoint() string {
	schema := "http"
	if c.MinIO.UseSSL {
		schema = "https"
	}
	return fmt.Sprintf("%s://%s:%s", schema, c.MinIO.Endpoint, c.MinIO.Port)
}

func (c *Config) GetCSRFAuthKey() []byte {
	csrfKey := getEnv("CSRF_AUTH_KEY", c.JWTSecret)
	if len(csrfKey) < 32 {
		return []byte(c.JWTSecret + "csrf-key-padding-to-32-chars")
	}
	return []byte(csrfKey[:32])
}
