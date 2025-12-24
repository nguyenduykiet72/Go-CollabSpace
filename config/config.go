package config

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	Database DBConfig     `yaml:"database"`
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

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		cfg = &Config{}

		_ = godotenv.Load()

		err := cleanenv.ReadConfig("./config/config.dev.yml", cfg)
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
