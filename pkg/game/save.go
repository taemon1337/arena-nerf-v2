package game

import (
  "os"
  "fmt"
  "path/filepath"
  "encoding/json"
)

func (ge *GameEngine) Logfile() string {
  return filepath.Join(ge.conf.Logdir, fmt.Sprintf("%s-%s.json", ge.CurrentGame.Mode(), ge.CurrentGame.Id()))
}

func (ge *GameEngine) Logstats() error {
  data, err := json.Marshal(ge.CurrentGameState)
  if err != nil {
    return err
  }

  return os.WriteFile(ge.Logfile(), data, os.ModePerm)
}

