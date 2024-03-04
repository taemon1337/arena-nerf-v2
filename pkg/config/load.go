package config

import (
  "os"
  "log"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "github.com/taemon1337/arena-nerf/pkg/common"
)

func (cfg *Config) HasConfig() bool {
  return cfg.ConfigFile != "" && common.FileExist(cfg.ConfigFile)
}

func LoadConfig(config_file string) (*Config, error) {
  logger := log.New(os.Stdout, "[config]: ", log.LstdFlags)
  cfg := NewConfig(logger)

  yamlFile, err := os.ReadFile(config_file)
  if err != nil {
    return cfg, err
  }

  err = yaml.Unmarshal(yamlFile, cfg)
  if err != nil {
    return cfg, err
  }

  return cfg, nil
}

func (cfg *Config) LoadConfig() error {
  yamlFile, err := os.ReadFile(cfg.ConfigFile)
  if err != nil {
    return err
  }

  err = yaml.Unmarshal(yamlFile, cfg)
  if err != nil {
    return err
  }

  return nil
}

func (cfg *Config) SaveConfig() error {
  yamlBytes, err := yaml.Marshal(cfg)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(cfg.ConfigFile, yamlBytes, 0640)
  return err
}
