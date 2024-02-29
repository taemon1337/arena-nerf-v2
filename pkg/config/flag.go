package config

import (
  "flag"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

func (c *Config) Flags() error {
  var tags []string

  flag.BoolVar(&c.EnableController, "enable-controller", c.EnableController, "enables the controller")
  flag.BoolVar(&c.EnableGameEngine, "enable-game-engine", c.EnableGameEngine, "enables the game engine on the controller")
  flag.BoolVar(&c.EnableNode, "enable-node", c.EnableNode, "enables the node")
  flag.BoolVar(&c.EnableSensor, "enable-sensor", c.EnableSensor, "enables sensors on this node")
  flag.BoolVar(&c.EnableSimulation, "enable-simulation", c.EnableSimulation, "enables game simulation")
  flag.BoolVar(&c.EnableConnector, "enable-connector", c.EnableConnector, "enables clustering with other nodes")

  flag.StringVar(&c.AgentConf.NodeName, "name", c.AgentConf.NodeName, "name of this node in the cluster")
  flag.StringVar(&c.AgentConf.BindAddr, "bind", c.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&c.AgentConf.AdvertiseAddr, "advertise", c.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&c.AgentConf.EncryptKey, "encrypt", c.AgentConf.EncryptKey, "encryption key")
  flag.BoolVar(&c.Coalesce, "coalesce", c.Coalesce, "enable to coalesce serf events sent to nodes")
  flag.Var((*AppendSliceValue)(&tags), "tag", "add tag to node with key=value")
  flag.Var((*AppendSliceValue)(&c.JoinAddrs), "join", "addresses to try to join automatically and repeatable until success")
  flag.Var((*AppendSliceValue)(&c.Nodes), "node", "add expected node by name, games will wait until all expected nodes are ready")
  flag.Var((*AppendSliceValue)(&c.Teams), "team", "add teams to be used in games")
  flag.IntVar(&c.Timeout, "timeout", c.Timeout, "number of seconds to wait to timeout nodes/connections/etc")
  flag.StringVar(&c.Logdir, "logdir", c.Logdir, "The directory to store game logs (which are served from the UI)")

  flag.Parse()

  parsedtags, err := UnmarshalTags(tags)
  if err != nil {
    return err
  }

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

  if err := c.Validate(); err != nil {
    return err
  }

  return nil
}
