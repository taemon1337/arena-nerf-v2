package config

import (
  "flag"
  "strings"
)

type AppendSliceValue []string

var _ flag.Value = new(AppendSliceValue)

func (s *AppendSliceValue) String() string {
  return strings.Join(*s, ",")
}

func (s *AppendSliceValue) Set(value string) error {
  if *s == nil {
    *s = make([]string, 0, 1)
  }

  *s = append(*s, value)
  return nil
}

