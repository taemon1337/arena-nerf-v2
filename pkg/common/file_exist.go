package common

import (
  "os"
  "log"
)

func FileExist(filepath string) bool {
  _, err := os.Stat(filepath)
  if err == nil {
    return true
  }
  if err == os.ErrNotExist {
    return false
  }

  log.Printf("WARNING: error checking file existance '%s': %s", filepath, err)
  return false
}
