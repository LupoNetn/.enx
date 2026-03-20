package utils

import (
	"os"
)

func ExtractKeyFromEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}