package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Queue    QueueConfig
	Server   ServerConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	External ExternalConfig
}

type AppConfig struct {
	Environment string
	Name        string
	Version     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type QueueConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	Secret string
}

// OAuthConfig OAuth configuration for all providers
type OAuthConfig struct {
	Google GoogleOAuthConfig `json:"google"`
}

type GoogleOAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	Scopes       string `json:"scopes"`
}

type ExternalConfig struct {
	AWS   AWSConfig
	Email EmailConfig
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3BucketName    string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Name:        getEnv("APP_NAME", "Weave"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "weave_user"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "weave_dev"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Queue: QueueConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "weave_user"),
			Password: getEnv("RABBITMQ_PASSWORD", ""),
			VHost:    getEnv("RABBITMQ_VHOST", "weave_vhost"),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
				Scopes:       getEnv("GOOGLE_SCOPES", "profile email"),
			},
		},
		External: ExternalConfig{
			AWS: AWSConfig{
				AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
				SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
				Region:          getEnv("AWS_REGION", "us-east-1"),
				S3BucketName:    getEnv("S3_BUCKET_NAME", ""),
			},
			Email: EmailConfig{
				SMTPHost:     getEnv("SMTP_HOST", ""),
				SMTPPort:     getEnv("SMTP_PORT", "587"),
				SMTPUsername: getEnv("SMTP_USERNAME", ""),
				SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}