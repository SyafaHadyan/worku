// Package env loads and parses environment variables from the .env file at the root directory
package env

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	LimiterMax                        int    `env:"LIMITER_MAX"`
	LimiterExpirationMinutes          int    `env:"LIMITER_EXPIRATION_MINUTES"`
	BodyLimit                         int    `env:"BODY_LIMIT_MB"`
	AccountRegistrationCodeDigitCount uint   `env:"ACCOUNT_REGISTRATION_CODE_DIGIT_COUNT"`
	PasswordChangeCodeDigitcount      uint   `env:"PASSWORD_CHANGE_CODE_DIGIT_COUNT"`
	PasswordChangeExpiryMinutes       int    `env:"PASSWORD_CHANGE_EXPIRY_MINUTES"`
	PasswordChangeCodeRetrySeconds    int    `env:"PASSWORD_CHANGE_CODE_RETRY_SECONDS"`
	AppPort                           uint   `env:"APP_PORT"`
	DBName                            string `env:"DB_NAME"`
	DBUsername                        string `env:"DB_USERNAME"`
	DBPassword                        string `env:"DB_PASSWORD"`
	DBHost                            string `env:"DB_HOST"`
	DBPort                            uint   `env:"DB_PORT"`
	RedisAddress                      string `env:"REDIS_ADDRESS"`
	RedisPort                         uint   `env:"REDIS_PORT"`
	RedisUsername                     string `env:"REDIS_USERNAME"`
	RedisPassword                     string `env:"REDIS_PASSWORD"`
	RedisDatabase                     int    `env:"REDIS_DATABASE"`
	RedisExpiration                   int    `env:"REDIS_EXPIRATION"`
	JWTSecretKey                      string `env:"JWT_SECRET_KEY"`
	JWTExpiredDays                    uint   `env:"JWT_EXPIRED_DAYS"`
	MidtransServerKey                 string `env:"MIDTRANS_SERVER_KEY"`
}

func New() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Panic("failed to load env")

		return nil
	}

	envParsed := new(Env)
	err = env.Parse(envParsed)
	if err != nil {
		log.Panic("failed to parse env")

		return nil
	}

	log.Println("loaded config")

	return envParsed
}
