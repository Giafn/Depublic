package configs

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `env:"ENV" envDefault:"development"`
	Port     string         `env:"PORT" envDefault:"8000"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	Encrypt  EncryptConfig  `envPrefix:"ENCRYPT_"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DBNAME" envDefault:"postgres"`
}

type JWTConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"secret"`
}

type RedisConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"6379"`
	Password string `env:"PASSWORD" envDefault:""`
}

type EncryptConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"secret"`
	Iv        string `env:"IV" envDefault:"iv"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("ERROR LOADING .ENV FILE")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.New("ERROR PARSING ENVIRONMENT VARIABLES")
	}

	return &cfg, nil
}
