package config

import (
	"fmt"
	"os"
)

func Env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func StorageAddress() string {
	address := Env("STORAGE_ADDRESS", "localhost")
	port := Env("STORAGE_PORT", "9000")
	return fmt.Sprintf("http://%s:%s", address, port)
}

func AccessLog() bool {
	return Env("ACCESS_LOG", "") == "true"
}
