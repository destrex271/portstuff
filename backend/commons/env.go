package commons

import (
	"syscall"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}

func SetEnvString(key, value string) error {
	if err := syscall.Setenv(key, value); err != nil {
		return err
	}
	return nil
}
