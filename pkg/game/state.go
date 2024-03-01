package game

import (
  "fmt"
  "sync"
  "time"
  "slices"
  "strings"
  "math/rand"
  "gopkg.in/yaml.v2"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameState struct {
  config            *GameConfig     `yaml:"config" json:"config"`
  Status            string          `yaml:"status" json:"status"`
  Teams             []string        `yaml:"teams" json:"teams"`
  Nodes             []string        `yaml:"nodes" json:"nodes"`
  Scoreboard        map[string]int  `yaml:"scoreboard" json:"scoreboard"`
  Nodeboard         map[string]int  `yaml:"nodeboard" json:"nodeboard"`
  Winner            string          `yaml:"winner" json:"winner"`
  Highscore         int             `yaml:"highscore" json:"highscore"`
  StartedAt         time.Time       `yaml:"StartedAt" json:"StartedAt"`
  GameDuration      time.Duration   `yaml:"GameDuration" json:"GameDuration"`
  EndedAt           time.Time       `yaml:"EndedAt" json:"EndedAt"`
  Timeline          []GameEvent     `yaml:"timeline" json:"timeline"`
  gamelock          *sync.Mutex     `yaml:"-" json:"-"`
}

func NewGameState(cfg *GameConfig) *GameState {
  return &GameState{
    config:         cfg,
    Status:         constants.GAME_STATUS_INIT,
    StartedAt:      time.Time{},
    EndedAt:        time.Time{},
    GameDuration:   0,
    Teams:          cfg.Cfg.Teams,
    Nodes:          cfg.Cfg.Nodes,
    Scoreboard:     map[string]int{},
    Nodeboard:      map[string]int{},
    Timeline:       make([]GameEvent, 0),
    gamelock:       &sync.Mutex{},
  }
}

func (gs *GameState) SetStatus(status string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.Status = status
}

func (gs *GameState) LogGameEvent(e GameEvent) {
  gs.Timeline = append(gs.Timeline, e)
}

func (gs *GameState) GameSetup() error {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.Status = constants.GAME_STATUS_RUNNING
  gs.StartedAt = time.Now()

  dur, err := time.ParseDuration(gs.config.GameLength)
  if err != nil {
    return err
  }

  gs.GameDuration = dur
  return nil
}

func (gs *GameState) EndGame() {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.Status = constants.GAME_STATUS_ENDED
  gs.EndedAt = time.Now()
}

func (gs *GameState) SetWinner(team string, score int) {
  if (score >= gs.Highscore) {
    gs.gamelock.Lock()
    defer gs.gamelock.Unlock()
    gs.Winner = team
    gs.Highscore = score
  } // we should probably return error in else case
}

func (gs *GameState) Running() bool {
  return gs.Status == constants.GAME_STATUS_RUNNING
}

func (gs *GameState) AddNode(node string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  if !slices.Contains(gs.Nodes, node) {
    gs.Nodes = append(gs.Nodes, node)
  }
}

func (gs *GameState) AddTeam(team string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  if !slices.Contains(gs.Teams, team) {
    gs.Teams = append(gs.Teams, team)
  }
}

func (gs *GameState) SetBoards(sb, nb map[string]int) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.Scoreboard = sb
  gs.Nodeboard = nb
}

func (gs *GameState) TeamList() string {
  return strings.Join(gs.Teams, constants.COMMA)
}

func (gs *GameState) RandomTeam() string {
  if len(gs.Teams) > 0 {
    return gs.Teams[rand.Intn(len(gs.Teams))]
  } else {
    return ""
  }
}

func (gs *GameState) RandomNode() string {
  if len(gs.Nodes) > 0 {
    return gs.Nodes[rand.Intn(len(gs.Nodes))]
  } else {
    return ""
  }
}

func (gs *GameState) Winners() []string {
  winningteams := []string{}

  for team, points := range gs.Scoreboard {
    if points >= gs.config.WinningScore {
      winningteams = append(winningteams, team)
    }
  }

  return winningteams
}

func (gs *GameState) ValidateNodes() error {
  l := len(gs.Nodes)
  t := len(gs.Teams)

  if l < gs.config.MinNodeCount {
    return constants.ERR_MIN_NODE_COUNT
  }

  if l > gs.config.MaxNodeCount {
    return constants.ERR_MAX_NODE_COUNT
  }

  if t < gs.config.MinTeamCount {
    return constants.ERR_MIN_TEAM_COUNT
  }

  if t > gs.config.MaxTeamCount {
    return constants.ERR_MAX_TEAM_COUNT
  }

  for _,  node := range gs.config.RequiredNodeNames {
    if !slices.Contains(gs.Nodes, node) {
      return constants.ERR_REQUIRED_NODE
    }
  }

  for _,  team := range gs.config.RequiredTeamNames {
    if !slices.Contains(gs.Teams, team) {
      return constants.ERR_REQUIRED_TEAM
    }
  }

  return nil
}

func (gs *GameState) TimeExpired() bool {
  if gs.Status == constants.GAME_STATUS_ENDED {
    return true
  }

  if time.Now().After(gs.StartedAt.Add(gs.GameDuration)) {
    return true
  }

  return false
}

func (gs *GameState) WinningScoreReached() bool {
  if gs.Status == constants.GAME_STATUS_ENDED {
    return true
  }

  if len(gs.Winners()) > 0 {
    return true
  }

  return false
}

func (gs *GameState) GameStatus() string {
  s := "###################################\n"
  s += fmt.Sprintf("Game Status: (%s)", gs.Status)
  s += fmt.Sprintf("Time Remaining: %s", gs.GameDuration - time.Since(gs.StartedAt))
  s += fmt.Sprintf("Scoreboard: \n%s", gs.Scoreboard)
  return s
}

func (gs *GameState) String() string {
  yamlBytes, err := yaml.Marshal(gs)
  if err != nil {
    return fmt.Sprintf("Error marshalling game state into yaml: %s", err)
  }
  return string(yamlBytes)
}
