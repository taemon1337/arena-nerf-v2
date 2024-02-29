package game

import (
  "sync"
  "time"
  "slices"
  "strings"
  "math/rand"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameState struct {
  config            *GameConfig     `yaml:"config" json:"config"`
  status            string          `yaml:"status" json:"status"`
  teams             []string        `yaml:"teams" json:"teams"`
  nodes             []string        `yaml:"nodes" json:"nodes"`
  scoreboard        map[string]int  `yaml:"scoreboard" json:"scoreboard"`
  nodeboard         map[string]int  `yaml:"nodeboard" json:"nodeboard"`
  winner            string          `yaml:"winner" json:"winner"`
  highscore         int             `yaml:"highscore" json:"highscore"`
  started_at        time.Time       `yaml:"started_at" json:"started_at"`
  game_duration     time.Duration   `yaml:"game_duration" json:"game_duration"`
  ended_at          time.Time       `yaml:"ended_at" json:"ended_at"`
  timeline          []GameEvent     `yaml:"timeline" json:"timeline"`
  gamelock          *sync.Mutex     `yaml:"-" json:"-"`
}

func NewGameState(cfg *GameConfig) *GameState {
  return &GameState{
    config:         cfg,
    status:         constants.GAME_STATUS_INIT,
    started_at:     time.Time{},
    ended_at:       time.Time{},
    game_duration:  0,
    teams:          cfg.Cfg.Teams,
    nodes:          cfg.Cfg.Nodes,
    scoreboard:     map[string]int{},
    nodeboard:      map[string]int{},
    timeline:       make([]GameEvent, 0),
    gamelock:       &sync.Mutex{},
  }
}

func (gs *GameState) SetStatus(status string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.status = status
}

func (gs *GameState) LogGameEvent(e GameEvent) {
  gs.timeline = append(gs.timeline, e)
}

func (gs *GameState) GameSetup() error {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.status = constants.GAME_STATUS_RUNNING
  gs.started_at = time.Now()

  dur, err := time.ParseDuration(gs.config.GameLength)
  if err != nil {
    return err
  }

  gs.game_duration = dur
  return nil
}

func (gs *GameState) EndGame() {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.status = constants.GAME_STATUS_ENDED
  gs.ended_at = time.Now()
}

func (gs *GameState) Winner() string {
  return gs.winner
}

func (gs *GameState) Highscore() int {
  return gs.highscore
}

func (gs *GameState) SetWinner(team string, score int) {
  if (score >= gs.highscore) {
    gs.gamelock.Lock()
    defer gs.gamelock.Unlock()
    gs.winner = team
    gs.highscore = score
  } // we should probably return error in else case
}

func (gs *GameState) Running() bool {
  return gs.status == constants.GAME_STATUS_RUNNING
}

func (gs *GameState) GameDuration() time.Duration {
  return gs.game_duration
}

func (gs *GameState) AddNode(node string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  if !slices.Contains(gs.nodes, node) {
    gs.nodes = append(gs.nodes, node)
  }
}

func (gs *GameState) AddTeam(team string) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  if !slices.Contains(gs.teams, team) {
    gs.teams = append(gs.teams, team)
  }
}

func (gs *GameState) Teams() []string {
  return gs.teams
}

func (gs *GameState) Nodes() []string {
  return gs.nodes
}

func (gs *GameState) SetBoards(sb, nb map[string]int) {
  gs.gamelock.Lock()
  defer gs.gamelock.Unlock()
  gs.scoreboard = sb
  gs.nodeboard = nb
}

func (gs *GameState) TeamList() string {
  return strings.Join(gs.teams, constants.COMMA)
}

func (gs *GameState) RandomTeam() string {
  if len(gs.teams) > 0 {
    return gs.teams[rand.Intn(len(gs.teams))]
  } else {
    return ""
  }
}

func (gs *GameState) RandomNode() string {
  if len(gs.nodes) > 0 {
    return gs.nodes[rand.Intn(len(gs.nodes))]
  } else {
    return ""
  }
}

func (gs *GameState) Winners() []string {
  winningteams := []string{}

  for team, points := range gs.scoreboard {
    if points >= gs.config.WinningScore {
      winningteams = append(winningteams, team)
    }
  }

  return winningteams
}

func (gs *GameState) ValidateNodes() error {
  l := len(gs.nodes)
  t := len(gs.teams)

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
    if !slices.Contains(gs.nodes, node) {
      return constants.ERR_REQUIRED_NODE
    }
  }

  for _,  team := range gs.config.RequiredTeamNames {
    if !slices.Contains(gs.teams, team) {
      return constants.ERR_REQUIRED_TEAM
    }
  }

  return nil
}

func (gs *GameState) Completed() bool {
  if gs.status == constants.GAME_STATUS_ENDED {
    return true
  }

  if len(gs.Winners()) > 0 {
    return true
  }

  if time.Now().After(gs.started_at.Add(gs.game_duration)) {
    return true
  }

  return false
}
