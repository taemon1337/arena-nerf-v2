package node

import (
  "sync"
  "strings"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type NodeState struct {
  Name          string          `yaml:"name" json:"name"`
  Status        string          `yaml:"status" json:"status"`
  Mode          string          `yaml:"mode" json:"mode"`
  Teams         []string        `yaml:"teams" json:"teams"`
  Colors        []string        `yaml:"colors" json:"colors"`
  Hits          map[string]int  `yaml:"hits" json:"hits"`
  nodelock      *sync.Mutex     `yaml:"-" json:"-"`
}

func NewNodeState(name string) *NodeState {
  return &NodeState{
    Name:         name,
    Status:       constants.GAME_STATUS_INIT,
    Mode:         "",
    Teams:        []string{},
    Colors:       []string{},
    Hits:         map[string]int{name: 0},
    nodelock:     &sync.Mutex{},
  }
}

func (ns *NodeState) SetTeams(teams string, enable_team_colors bool) {
  ns.nodelock.Lock()
  defer ns.nodelock.Unlock()
  ns.Teams = strings.Split(teams, constants.COMMA)
  if enable_team_colors {
    ns.Colors = ns.Teams
  }
}

func (ns *NodeState) AddTeamHit(team string, count int) {
  ns.AddNodeHit(constants.NONE_SENSOR_ID, team, count)
}

func (ns *NodeState) AddNodeHit(sensorid, sensorcolor string, hitcount int) {
  ns.nodelock.Lock()
  defer ns.nodelock.Unlock()

  // TODO: we should check if EnableTeamColors is set
  if _, ok := ns.Hits[sensorid]; !ok {
    ns.Hits[sensorid] = 0 // initialize
  }

  if _, ok := ns.Hits[sensorcolor]; !ok {
    ns.Hits[sensorcolor] = 0 // initialize
  }

  ns.Hits[ns.Name] += hitcount // total node hits
  ns.Hits[sensorid] += hitcount // total sensor hits
  ns.Hits[sensorcolor] += hitcount // total team/color hits
}
