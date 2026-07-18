package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DBConfig          `yaml:"database"`
	JWT         JWTConfig         `yaml:"jwt"`
	Redis       RedisConfig       `yaml:"redis"`
	AWS         AWSConfig         `yaml:"aws"`
	ResendEmail ResendEmailConfig `yaml:"resendEmail"`
	//SMTP        SMTPConfig        `yaml:"smtp"`
	OAuthGoogleConfig OAuthGoogleConfig `yaml:"oauthGoogle"`
}

type ServerConfig struct {
	Port                       int      `yaml:"port" env:"PORT" env-default:"8080" validate:"required"`
	Mode                       string   `yaml:"mode" env:"MODE" env-default:"development" validate:"required"`
	AllowedOrigins             []string `yaml:"allowedOrigins" env:"ALLOWED_ORIGINS" env-separator:"," env-default:"http://localhost:3000,http://localhost:3001"`
	FrontendReturnURLWhitelist []string `yaml:"frontendReturnURLWhitelist" env:"FE_RETURN_URL_WHITELIST" env-separator:"," env-default:"http://localhost:3000/reset-password,http://localhost:3001/reset-password"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" validate:"required"`
	Port     int    `yaml:"port" env:"DB_PORT" validate:"required"`
	User     string `yaml:"username" env:"DB_USER" validate:"required"`
	Password string `yaml:"password" env:"DB_PASSWORD" validate:"required"`
	DBName   string `yaml:"dbname" env:"DB_NAME" validate:"required"`
	Type     string `yaml:"type" env:"DB_TYPE" env-default:"postgres" validate:"required,oneof=postgres mysql"`
}

type OAuthGoogleConfig struct {
	ClientID     string `yaml:"client_id" env:"GOOGLE_CLIENT_ID" validate:"required"`
	ClientSecret string `yaml:"client_secret" env:"GOOGLE_CLIENT_SECRET" validate:"required"`
	RedirectURL  string `yaml:"redirect_url" env:"GOOGLE_REDIRECT_URL" validate:"required"`
}

type JWTConfig struct {
	AccessTokenSecret    string        `yaml:"accessTokenSecret" env:"JWT_ACCESS_SECRET" validate:"required"`
	RefreshTokenSecret   string        `yaml:"refreshTokenSecret" env:"JWT_REFRESH_SECRET" validate:"required"`
	AccessTokenDuration  time.Duration `yaml:"accessTokenDuration" env:"JWT_ACCESS_DURATION" env-default:"1h"`
	RefreshTokenDuration time.Duration `yaml:"refreshTokenDuration" env:"JWT_REFRESH_DURATION" env-default:"168h"` // 7 days
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" validate:"required"`
	Port     int    `yaml:"port" env:"REDIS_PORT" validate:"required"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" validate:"required"`
}

type AWSConfig struct {
	Region          string `yaml:"region" env:"AWS_REGION" validate:"required"`
	AccessKeyID     string `yaml:"accessKeyId" env:"AWS_ACCESS_KEY_ID" validate:"required"`
	SecretAccessKey string `yaml:"secretAccessKey" env:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	BucketName      string `yaml:"bucketName" env:"AWS_BUCKET_NAME" validate:"required"`
}

type ResendEmailConfig struct {
	ResendAPIKey string `yaml:"resendApiKey" env:"RESEND_API_KEY" validate:"required"`
	FromEmail    string `yaml:"fromEmail" env:"RESEND_FROM_EMAIL" validate:"required"`
}

//type SMTPConfig struct {
//	Host     string `yaml:"host" env:"SMTP_HOST" env-default:"smtp.gmail.com"`
//	Port     int    `yaml:"port" env:"SMTP_PORT" env-default:"587"`
//	Username string `yaml:"username" env:"SMTP_USERNAME" validate:"required"`
//	Password string `yaml:"password" env:"SMTP_PASSWORD" validate:"required"` // Dùng App Password nếu là Gmail
//	From     string `yaml:"from" env:"SMTP_FROM" validate:"required"`
//}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		cfg = &Config{}

		_ = godotenv.Load()

		path := resolveConfigPath()
		if path != "" {
			if readErr := cleanenv.ReadConfig(path, cfg); readErr != nil {
				err = fmt.Errorf("failed to read config %q: %w", path, readErr)
				return
			}
		} else {
			if envErr := cleanenv.ReadEnv(cfg); envErr != nil {
				err = fmt.Errorf("failed to read env config: %w", envErr)
				return
			}
		}

		validate := validator.New()
		if validateErr := validate.Struct(cfg); validateErr != nil {
			err = fmt.Errorf("config validation failed: %w", validateErr)
			return
		}
	})

	return cfg, err
}

// resolveConfigPath picks the YAML file to read. Returns "" when no file is
// available, signalling that the caller should fall back to env-only loading.
func resolveConfigPath() string {
	if explicit := strings.TrimSpace(os.Getenv("CONFIG_PATH")); explicit != "" {
		if fileExists(explicit) {
			return explicit
		}
		fmt.Fprintf(os.Stderr, "CONFIG_PATH=%q not found, falling back to environment variables\n", explicit)
		return ""
	}

	mode := strings.TrimSpace(os.Getenv("MODE"))
	if mode == "" {
		mode = "development"
	}
	candidate := filepath.Join("config", fmt.Sprintf("config.%s.yml", mode))
	if fileExists(candidate) {
		return candidate
	}

	return ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		fmt.Fprintf(os.Stderr, "config stat %q: %v\n", path, err)
		return true
	}
	return !info.IsDir()
}
