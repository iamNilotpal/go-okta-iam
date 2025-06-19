package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	JWT      *JWTConfig
	Okta     *OktaConfig
	Server   *ServerConfig
	Frontend *FrontendConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type OktaConfig struct {
	Domain       string
	ClientID     string
	ClientSecret string
	APIToken     string
	Issuer       string
	Audience     string
	RedirectURI  string
}

type FrontendConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	oktaDomain := os.Getenv("OKTA_DOMAIN")

	config := &Config{
		JWT:      &JWTConfig{Secret: os.Getenv("JWT_SECRET")},
		Frontend: &FrontendConfig{URL: os.Getenv("FRONTEND_URL")},
		Server: &ServerConfig{
			Port:         getEnvOrDefault("PORT", "8080"),
			ReadTimeout:  getDurationOrDefault("READ_TIMEOUT", "10s"),
			WriteTimeout: getDurationOrDefault("WRITE_TIMEOUT", "10s"),
			IdleTimeout:  getDurationOrDefault("IDLE_TIMEOUT", "120s"),
		},
		Okta: &OktaConfig{
			Domain:       oktaDomain,
			Issuer:       os.Getenv("OKTA_ISSUER"),
			Audience:     os.Getenv("OKTA_AUDIENCE"),
			ClientID:     os.Getenv("OKTA_CLIENT_ID"),
			APIToken:     os.Getenv("OKTA_API_TOKEN"),
			RedirectURI:  os.Getenv("OKTA_REDIRECT_URI"),
			ClientSecret: os.Getenv("OKTA_CLIENT_SECRET"),
		},
	}

	return config, nil
}

func (c *Config) GetOktaIssuer() string {
	return c.Okta.Issuer
}

func (c *Config) GetOktaJWKSURL() string {
	return fmt.Sprintf("%s/v1/keys", c.Okta.Issuer)
}

func (c *Config) GetOktaAuthURL() string {
	return fmt.Sprintf("https://%s/oauth2/default/v1/authorize", c.Okta.Domain)
}

func (c *Config) GetOktaTokenURL() string {
	return fmt.Sprintf("https://%s/oauth2/default/v1/token", c.Okta.Domain)
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
