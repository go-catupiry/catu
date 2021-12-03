package configuration

import (
	"os"
	"strconv"
)

// Configuration object with usefull methods
type Cfg struct {
}

func (c Cfg) Init() error {
	// TODO!
	return nil
}

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
