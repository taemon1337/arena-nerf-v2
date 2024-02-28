package game

type GameConfig struct {
  GameLength                string        `yaml:"game_length" json:"game_length"`
  WinningScore              int           `yaml:"winning_score" json:"winning_score"`
  MinNodeCount              int           `yaml:"max_node_count" json:"max_node_count"`
  MaxNodeCount              int           `yaml:"min_node_count" json:"min_node_count"`
  MinTeamCount              int           `yaml:"min_team_count" json:"min_team_count"`
  MaxTeamCount              int           `yaml:"max_team_count" json:"max_team_count"`
  RequiredNodeNames         []string      `yaml:"required_node_names" json:"required_node_names"`
  RequiredTeamNames         []string      `yaml:"required_team_names" json:"required_team_names"`
}

func NewGameConfig() *GameConfig {
  return &GameConfig{
    GameLength:         "3m",
    WinningScore:       100,
    MinNodeCount:       3,
    MaxNodeCount:       100,
    MinTeamCount:       3,
    MaxTeamCount:       100,
    RequiredNodeNames:  []string{},
    RequiredTeamNames:  []string{},
  }
}
