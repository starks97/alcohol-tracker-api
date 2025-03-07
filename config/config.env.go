package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DatabaseUrl            string
	ClientOrigin           string
	Domain                 string
	PostgresUser           string
	PostgresPassword       string
	PostgresDB             string
	PostgresHost           string
	PostgresPort           string
	PostgresSSLMode        string
	GoogleLoginConfig      oauth2.Config
	GoogleClientID         string
	GoogleClientSecret     string
	RedisURL               string
	RedisPassword          string
	RedisAddress           string
	AccessTokenPrivateKey  string
	AccessTokenPublicKey   string
	AccessTokenExpiredIn   string
	AccessTokenMaxAge      int64
	RefreshTokenPrivateKey string
	RefreshTokenPublicKey  string
	RefreshTokenMaxAge     int64
	RefreshTokenExpiredIn  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}
	accessTokenMaxAge, err := parseDuration(getEnv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_MAXAGE: %v", err)
	}

	refreshTokenMaxAge, err := parseDuration(getEnv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_MAXAGE: %v", err)
	}

	config := &Config{
		DatabaseUrl:            getEnv("DATABASE_URL"),
		ClientOrigin:           getEnv("CLIENT_ORIGIN"),
		Domain:                 getEnv("DOMAIN"),
		PostgresUser:           getEnv("POSTGRES_USER"),
		PostgresPassword:       getEnv("POSTGRES_PASSWORD"),
		PostgresDB:             getEnv("POSTGRES_DB"),
		PostgresHost:           getEnv("POSTGRES_HOST"),
		PostgresPort:           getEnv("POSTGRES_PORT"),
		PostgresSSLMode:        getEnv("POSTGRES_SSLMODE"),
		GoogleClientID:         getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:     getEnv("GOOGLE_CLIENT_SECRET"),
		RedisURL:               getEnv("REDIS_URL"),
		RedisAddress:           getEnv("REDIS_ADDRESS"),
		RedisPassword:          getEnv("REDIS_PASSWORD"),
		AccessTokenPrivateKey:  getEnv("ACCESS_TOKEN_PRIVATE_KEY"),
		AccessTokenPublicKey:   getEnv("ACCESS_TOKEN_PUBLIC_KEY"),
		AccessTokenExpiredIn:   getEnv("ACCESS_TOKEN_EXPIRED_IN"),
		AccessTokenMaxAge:      accessTokenMaxAge,
		RefreshTokenPrivateKey: getEnv("REFRESH_TOKEN_PRIVATE_KEY"),
		RefreshTokenPublicKey:  getEnv("REFRESH_TOKEN_PUBLIC_KEY"),
		RefreshTokenMaxAge:     refreshTokenMaxAge,
		RefreshTokenExpiredIn:  getEnv("REFRESH_TOKEN_EXPIRED_IN"),
	}

	// Initialize OAuth2 configuration
	config.GoogleLoginConfig = oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  "http://localhost:8080/auth/google_callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return config, nil
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s must be set", key)
	}

	return value
}

func parseDuration(durationStr string) (int64, error) {
	var numericPart string

	for _, char := range durationStr {
		if unicode.IsDigit(char) {
			numericPart += string(char)
		} else {
			break
		}
	}

	if numericPart == "" {
		return 0, fmt.Errorf("invalid duration string")
	}

	duration, err := strconv.ParseInt(numericPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric part: %v", err)
	}
	return duration, nil
}
