package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(pathname string) error {
	err := godotenv.Load(pathname)
	if err != nil {
		return err
	}
	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}