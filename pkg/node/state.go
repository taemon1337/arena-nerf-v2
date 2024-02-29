package node

import (
  "sync"
  "strings"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type NodeState struct {
  name          string          `yaml:"name" json:"name"`
  status        string          `yaml:"status" json:"status"`
  mode          string          `yaml:"mode" json:"mode"`
  teams         []string        `yaml:"teams" json:"teams"`
  hits          map[string]int  `yaml:"hits" json:"hits"`
  nodelock      *sync.Mutex     `yaml:"-" json:"-"`
}

func NewNodeState(name string) *NodeState {
  return &NodeState{
    name:         name,
    status:       constants.GAME_STATUS_INIT,
    mode:         "",
    teams:        []string{},
    hits:         map[string]int{name: 0},
    nodelock:     &sync.Mutex{},
  }
}

func (ns *NodeState) Status() string {
  return ns.status
}

func (ns *NodeState) SetStatus(status string) {
  ns.status = status
}

func (ns *NodeState) Hits() map[string]int {
  return ns.hits
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

func (ns *NodeState) AddTeamHit(team string, count int) {
  if _, ok := ns.hits[team]; ok {
    ns.hits[team] += count
  } else {
    ns.hits[team] = count
  }
}

func (ns *NodeState) AddNodeHit(count int) {
  ns.hits[ns.name] += count
}
