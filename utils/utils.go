package utils

import (
	"os"
	"strings"
)

func GetEnv(key string, d string) string {
	val := os.Getenv(key)
	if val == "" {
		return d
	}
	return val
}

func GetEnvBool(key string, d bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return d
	}
	if strings.ToLower(val) == "true" {
		return true
	}

	return false
}
