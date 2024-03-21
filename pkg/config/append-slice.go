package config

import (
  "flag"
  "slices"
  "strings"
)

type AppendSliceValue []string

var _ flag.Value = new(AppendSliceValue)

func (s *AppendSliceValue) String() string {
  return strings.Join(*s, ",")
}

func (s *AppendSliceValue) Set(value string) error {
  if *s == nil {
    *s = []string{}
  }

  if !slices.Contains(*s, value) {
    *s = append(*s, value)
  }
  return nil
}

