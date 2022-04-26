package configuration

import (
	"os"
	"strconv"
)

// Default interface to be used on others modules
type ConfigurationInterface interface {
	Init() error
	Get(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64

	GetF(key, fallback string) string
	GetBoolF(key string, fallback bool) bool
	GetIntF(key string, fallback int) int
	GetInt64F(key string, fallback int64) int64
}

// Build and get a new Cfg object
func NewCfg() ConfigurationInterface {
	c := Cfg{}

	return c
}

// Get environment variable value as string
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get environment variable value as boolean
func GetBoolEnv(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		boolV, err := strconv.ParseBool(value)
		if err == nil {
			return boolV
		}

		return false
	}

	return fallback
}

// Get environment variable value as int
func GetIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(value)
		if err == nil {
			return v
		}

		return 0
	}

	return fallback
}

// Get environment variable value as int64
func GetInt64Env(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return v
		}

		return 0
	}

	return fallback
}
