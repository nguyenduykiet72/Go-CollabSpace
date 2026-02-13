package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	Database DBConfig     `yaml:"database"`
	JWT      JWTConfig    `yaml:"jwt"`
	Redis    RedisConfig  `yaml:"redis"`
}

type ServerConfig struct {
	Port int    `yaml:"port" env:"PORT" env-default:"8080" validate:"required"`
	Mode string `yaml:"mode" env:"MODE" env-default:"development" validate:"required"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" validate:"required"`
	Port     int    `yaml:"port" env:"DB_PORT" validate:"required"`
	User     string `yaml:"username" env:"DB_USER" validate:"required"`
	Password string `yaml:"password" env:"DB_PASSWORD" validate:"required"`
	DBName   string `yaml:"dbname" env:"DB_NAME" validate:"required"`
	Type     string `yaml:"type" env:"DB_TYPE" env-default:"postgres" validate:"required,oneof=postgres mysql"`
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

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		cfg = &Config{}

		_ = godotenv.Load()

		err = cleanenv.ReadConfig("./config/config.dev.yml", cfg)
		if err != nil {
			fmt.Println("Config file not found, trying environment variables:", err)
			if envErr := cleanenv.ReadEnv(cfg); envErr != nil {
				err = envErr
				return
			}
		}

		validate := validator.New()
		if validateErr := validate.Struct(cfg); validateErr != nil {
			err = fmt.Errorf("config validation failed :%w", validateErr)
			return
		}
	})

	return cfg, err
}
