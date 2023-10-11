package env

import (
	"errors"
	"os"
)

func getEnv(key string) (string, error) {
  env := os.Getenv(key)
  if env == "" {
    return "", errors.New("Environment variable" + key + " is not provided.")
  }

  return env, nil
}

func GetSecretKey() (string, error) {
  return getEnv("SECRET_KEY")
}

func GetDbDns() (string, error) {
  return getEnv("POSTGRES_URL")
}
