package config

import "time"

type Config struct {
	Profiles map[string]Profile
	Current  string
}

type Profile struct {
	BaseURL     string
	Headers     map[string]string
	Auth        AuthConfig
	RetryConfig RetryConfig
}

type AuthConfig struct {
	Type     string // basic, oauth2, jwt, apikey
	Token    string
	Username string
	Password string
}

type RetryConfig struct {
	MaxRetries      int
	BackoffDuration time.Duration
	MaxDuration     time.Duration
}
