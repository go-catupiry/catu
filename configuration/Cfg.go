package configuration

import (
	"os"
	"strconv"
)

// Configuration object with usefull methods
type Cfg struct {
}

func (c Cfg) Init() error {
	return nil
}

// Get - Get environment variable value or "" if not exists
func (c Cfg) Get(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

// GetBoolEnv - Get an boolean env var. This returns false to invalid values
func (c Cfg) GetBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		boolV, err := strconv.ParseBool(value)
		if err == nil {
			return boolV
		}

		return false
	}

	return false
}

func (c Cfg) GetInt(key string) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(value)
		if err == nil {
			return v
		}

		return 0
	}

	return 0
}

func (c Cfg) GetInt64(key string) int64 {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return v
		}

		return 0
	}

	return 0
}

func (c Cfg) GetF(key, fallback string) string {
	return GetEnv(key, fallback)
}

// GetBoolEnv - Get an boolean environment var with default value. This returns false to invalid values
func (c Cfg) GetBoolF(key string, fallback bool) bool {
	return GetBoolEnv(key, fallback)
}

func (c Cfg) GetIntF(key string, fallback int) int {
	return GetIntEnv(key, fallback)
}

func (c Cfg) GetInt64F(key string, fallback int64) int64 {
	return GetInt64Env(key, fallback)
}
