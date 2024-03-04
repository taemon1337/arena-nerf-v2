package config

import (
  "flag"
  "slices"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

func (c *Config) Flags() error {
  var tags []string

  flag.StringVar(&c.ConfigFile, "config-file", c.ConfigFile, "path to read/write config yaml to")
  flag.BoolVar(&c.EnableController, "enable-controller", c.EnableController, "enables the controller")
  flag.BoolVar(&c.EnableGameEngine, "enable-game-engine", c.EnableGameEngine, "enables the game engine on the controller")
  flag.BoolVar(&c.EnableNode, "enable-node", c.EnableNode, "enables the node")
  flag.BoolVar(&c.EnableSensors, "enable-sensors", c.EnableSensors, "enables sensors on this node")
  flag.BoolVar(&c.EnableSimulation, "enable-simulation", c.EnableSimulation, "enables game simulation")
  flag.BoolVar(&c.EnableConnector, "enable-connector", c.EnableConnector, "enables clustering with other nodes")
  flag.BoolVar(&c.EnableTeamColors, "enable-team-colors", c.EnableTeamColors, "if set, all teams are also used as sensor led colors")

  flag.StringVar(&c.NodeName, "name", c.NodeName, "name of this node in the cluster")
  flag.StringVar(&c.AgentConf.BindAddr, "bind", c.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&c.AgentConf.AdvertiseAddr, "advertise", c.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&c.AgentConf.EncryptKey, "encrypt", c.AgentConf.EncryptKey, "encryption key")
  flag.BoolVar(&c.Coalesce, "coalesce", c.Coalesce, "enable to coalesce serf events sent to nodes")
  flag.Var((*AppendSliceValue)(&tags), "tag", "add tag to node with key=value")
  flag.Var((*AppendSliceValue)(&c.JoinAddrs), "join", "addresses to try to join automatically and repeatable until success")
  flag.Var((*AppendSliceValue)(&c.Nodes), "node", "add expected node by name, games will wait until all expected nodes are ready")
  flag.Var((*AppendSliceValue)(&c.Teams), "team", "add teams to be used in games")
  flag.Var((*AppendSliceValue)(&c.Colors), "color", "add color to available colors for LEDs")
  flag.IntVar(&c.Timeout, "timeout", c.Timeout, "number of seconds to wait to timeout nodes/connections/etc")
  flag.StringVar(&c.Logdir, "logdir", c.Logdir, "The directory to store game logs (which are served from the UI)")

  // -sensor 1:orangepi:gpiochip0:73:3
  flag.Var(c.SensorsConf, "sensor", "Add a sensor in the form of -sensor one:orangepi:gpiochip0:73:13, <1-4>:<device>:<gpiochip>:<hitpin>:<ledpin>")

  flag.Parse()

  parsedtags, err := UnmarshalTags(tags)
  if err != nil {
    return err
  }

  if c.HasConfig() {
    c.Printf("reading config from %s", c.ConfigFile)
    if err := c.LoadConfig(); err != nil {
      c.Printf("error reading config file: %s", c.ConfigFile)
      return err
    }
  }

  c.AgentConf.NodeName = c.NodeName
  c.AgentConf.Tags = parsedtags

  if c.EnableNode {
    c.Nodes = append(c.Nodes, c.AgentConf.NodeName)
  }

  if c.EnableController {
    c.AgentConf.Tags[constants.TAG_CTRL] = constants.TAG_TRUE
  }

  if c.EnableNode {
    c.AgentConf.Tags[constants.TAG_NODE] = constants.TAG_TRUE
  }

  if c.EnableTeamColors {
    for _, team := range c.Teams {
      if !slices.Contains(c.Colors, team) {
        c.Colors = append(c.Colors, team)
      }
    }
  }

  if err := c.Validate(); err != nil {
    return err
  }

  if c.ConfigFile != "" {
    c.Printf("saving config to %s", c.ConfigFile)
    if err := c.SaveConfig(); err != nil {
      c.Printf("error saving config file: %s", err)
      return err
    }
  }

  return nil
}
