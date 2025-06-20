package config

import (
	"os"
	"time"
)

type Config struct {
	Okta   *OktaConfig
	Server *ServerConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type OktaConfig struct {
	Domain   string
	APIToken string
	Issuer   string
	Audience string
}

type FrontendConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
}

func Load() (*Config, error) {
	config := &Config{
		Server: &ServerConfig{
			Port:         getEnvOrDefault("PORT", "8080"),
			ReadTimeout:  getDurationOrDefault("READ_TIMEOUT", "10s"),
			WriteTimeout: getDurationOrDefault("WRITE_TIMEOUT", "10s"),
			IdleTimeout:  getDurationOrDefault("IDLE_TIMEOUT", "120s"),
		},
		Okta: &OktaConfig{
			Domain:   os.Getenv("OKTA_DOMAIN"),
			Issuer:   os.Getenv("OKTA_ISSUER"),
			Audience: os.Getenv("OKTA_AUDIENCE"),
			APIToken: os.Getenv("OKTA_API_TOKEN"),
		},
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationOrDefault(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
