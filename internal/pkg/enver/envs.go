package enver

import (
  "fmt"
  "os"
)

func Env(key string) (string, error) {
  if val := os.Getenv(key); val != "" {
    return val, nil
  }
  return "", fmt.Errorf("%s value not found", key)
}
