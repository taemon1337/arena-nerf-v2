package config

import (
  "fmt"
  "strings"
)

func UnmarshalTags(tags []string) (map[string]string, error) {
  result := make(map[string]string)
  for _, tag := range tags {
    parts := strings.SplitN(tag, "=", 2)
    if len(parts) != 2 || len(parts[0]) == 0 {
      return nil, fmt.Errorf("Invalid tag: '%s'", tag)
    }
    result[parts[0]] = parts[1]
  }
  return result, nil
}

