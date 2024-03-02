package config

import (
  "os"
  "fmt"
  "log"
  "strings"

  "gopkg.in/yaml.v2"
  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"

  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type Config struct {
  EnableController        bool        `yaml:"enable_controller" json:"enable_controller"`
  EnableGameEngine        bool        `yaml:"enable_game_engine" json:"enable_game_engine"`
  EnableNode              bool        `yaml:"enable_node" json:"enable_node"`
  EnableSensors           bool        `yaml:"enable_sensors" json:"enable_sensors"`
  EnableSimulation        bool        `yaml:"enable_simulation" json:"enable_simulation"`
  EnableConnector         bool        `yaml:"enable_connector" json:"enable_connector"`
  EnableTeamColors        bool        `yaml:"enable_team_colors" json:"enable_team_colors"`
  Teams                   []string    `yaml:"teams" json:"teams"`
  Nodes                   []string    `yaml:"nodes" json:"nodes"`
  Colors                  []string    `yaml:"colors" json:"colors"`

  // serf config
  AgentConf               *agent.Config   `yaml:"-" json:"-"`
  SerfConf                *serf.Config    `yaml:"-" json:"-"`
  SensorsConf             *SensorsConfig  `yaml:"sensors" json:"sensors"`
  Coalesce                bool            `yaml:"coalesce" json:"coalesce"`
  JoinAddrs               []string        `yaml:"join_addrs" json:"join_addrs"`

  // game config
  WinningScore            int             `yaml:"winning_score" json:"winning_score"`
  GameLength              string          `yaml:"game_length" json:"game_length"`

  Timeout                 int             `yaml:"timeout" json:"timeout"`
  Logdir                  string          `yaml:"logdir" json:"logdir"`
  *log.Logger                             `yaml:"-" json:"-"`
}

func NewConfig(logger *log.Logger) *Config {
  ac := agent.DefaultConfig()
  sc := serf.DefaultConfig()
  joinaddrs := Getenv("SERF_JOIN_ADDRS", "127.0.0.1")
  ac.NodeName = Getenv("SERF_NAME", GetHostname())
  ac.BindAddr = Getenv("SERF_BIND_ADDR", "0.0.0.0")
  ac.AdvertiseAddr = Getenv("SERF_ADVERTISE_ADDR", "")
  ac.EncryptKey = Getenv("SERF_ENCRYPT_KEY", "")
  ac.LogLevel = Getenv("SERF_LOG_LEVEL", "err")

  return &Config{
    EnableController:   false,
    EnableGameEngine:   false,
    EnableNode:         false,
    EnableSensors:      false,
    EnableSimulation:   false,
    EnableConnector:    false,
    EnableTeamColors:   false,
    Teams:              []string{constants.BLUE_TEAM, constants.RED_TEAM, constants.YELLOW_TEAM, constants.GREEN_TEAM},
    Nodes:              []string{},
    Colors:             []string{},
    AgentConf:          ac,
    SerfConf:           sc,
    SensorsConf:        NewSensorsConfig(),
    Coalesce:           false,
    JoinAddrs:          strings.Split(joinaddrs, ","),
    WinningScore:       10,
    GameLength:         "3m",
    Timeout:            10, // 10 second timeouts
    Logdir:             "./logs",
    Logger:             log.New(logger.Writer(), "[CONFIG]: ", logger.Flags()),
  }
}

func (c *Config) Validate() error {
  ac := c.AgentConf
  sc := c.SerfConf

  var bindIP string
  var bindPort int
  var advertIP string
  var advertPort int
  var err error

  if ac.BindAddr != "" {
    bindIP, bindPort, err = ac.AddrParts(ac.BindAddr)
    if err != nil {
      return err
    }
  }

  if ac.AdvertiseAddr != "" {
    advertIP, advertPort, err = ac.AddrParts(ac.AdvertiseAddr)
    if err != nil {
      return err
    }
  }

  encryptKey, err := ac.EncryptBytes()
  if err != nil {
    return err
  }

  // https://github.com/hashicorp/serf/blob/master/cmd/serf/command/agent/command.go#L320
  sc.Tags = ac.Tags
  sc.NodeName = ac.NodeName
  sc.MemberlistConfig.BindAddr = bindIP
  sc.MemberlistConfig.BindPort = bindPort
  sc.MemberlistConfig.AdvertiseAddr = advertIP
  sc.MemberlistConfig.AdvertisePort = advertPort
  sc.MemberlistConfig.SecretKey = encryptKey
  sc.ProtocolVersion = uint8(ac.Protocol)
  sc.SnapshotPath = ac.SnapshotPath
  sc.MemberlistConfig.EnableCompression = ac.EnableCompression
  sc.QuerySizeLimit = ac.QuerySizeLimit
  sc.UserEventSizeLimit = ac.UserEventSizeLimit
  sc.EnableNameConflictResolution = !ac.DisableNameResolution
  sc.RejoinAfterLeave = ac.RejoinAfterLeave

  return nil
}

func Getenv(key, val string) string {
  a, exists := os.LookupEnv(key)
  if a != "" && exists {
    return a
  }
  return val // default
}

func GetHostname() string {
  hostname, err := os.Hostname()
  if err != nil {
    log.Printf("cannot get hostname - %s", err)
    return ""
  }

  return hostname
}

func (cfg *Config) String() string {
  yamlBytes, err := yaml.Marshal(cfg)
  if err != nil {
    return fmt.Sprintf("Error marshalling cfg into yaml: %s", err)
  }
  return string(yamlBytes)
}
