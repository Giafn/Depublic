package configs

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `env:"ENV" envDefault:"development"`
	Host     string         `env:"HOST" envDefault:"localhost"`
	Port     string         `env:"PORT" envDefault:"8000"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	Encrypt  EncryptConfig  `envPrefix:"ENCRYPT_"`
	SMTP     SMTPConfig     `envPrefix:"SMTP_"`
	Midtrans MidtransConfig `envPrefix:"MIDTRANS_"`
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
	ExpiresAt int    `env:"EXPIRES_AT" envDefault:"24"`
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

type SMTPConfig struct {
	Host string `env:"HOST" envDefault:"smtp.gmail.com"`
	Port string `env:"PORT" envDefault:"587"`
	User string `env:"USER" envDefault:""`
	Pass string `env:"PASS" envDefault:""`
	From string `env:"FROM" envDefault:"depublic@gmail.com"`
}

type MidtransConfig struct {
	ClientKey    string `env:"CLIENT_KEY" envDefault:""`
	ServerKey    string `env:"SERVER_KEY" envDefault:""`
	IsProduction bool   `env:"IS_PRODUCTION" envDefault:"false"`
	URL          string `env:"URL" envDefault:"https://api.sandbox.midtrans.com/snap/v1"`
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
