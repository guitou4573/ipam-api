package utils

import (
	"os"
)

func EnvGet (key string, defaultVal string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }

  return defaultVal
}
