package node

import (
  "sync"
  "strings"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type NodeState struct {
  name          string          `yaml:"name" json:"name"`
  mode          string          `yaml:"mode" json:"mode"`
  teams         []string        `yaml:"teams" json:"teams"`
  hits          map[string]int  `yaml:"hits" json:"hits"`
  nodelock      *sync.Mutex     `yaml:"-" json:"-"`
}

func NewNodeState(name string) *NodeState {
  return &NodeState{
    name:         name,
    mode:         "",
    teams:        []string{},
    hits:         map[string]int{},
    nodelock:     &sync.Mutex{},
  }
}

func (ns *NodeState) SetName(name string) {
  ns.nodelock.Lock()
  defer ns.nodelock.Unlock()
  ns.name = name
}

func (ns *NodeState) SetMode(mode string) {
  ns.nodelock.Lock()
  defer ns.nodelock.Unlock()
  ns.mode = mode
}

func (ns *NodeState) SetTeams(teams string) {
  ns.nodelock.Lock()
  defer ns.nodelock.Unlock()
  ns.teams = strings.Split(teams, constants.COMMA)
}

func (ns *NodeState) GetMode() string {
  return ns.mode
}
