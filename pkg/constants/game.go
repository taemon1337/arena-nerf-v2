package constants

var (
  // game modes
  GAME_MODE = "game:mode" // set game mode
  GAME_MODE_SIMULATION = "simulation"

  // game actions
  GAME_ACTION_BEGIN = "game:begin"
  GAME_ACTION_END = "game:end"
  GAME_ACTION_PAUSE = "game:pause"
  GAME_ACTION_RESET = "game:reset"

  // game states
  GAME_STATUS_INIT = "game:init"
  GAME_STATUS_RUNNING = "game:running"
  GAME_STATUS_ENDED = "game:over"
  GAME_STATUS_FAILED = "game:failed"
  GAME_TEAMS = "game:teams"
  GAME_WINNER = "game:winner"
  GAME_ERROR = "game:error"

  NODE_SCOREBOARD = "node:scoreboard"

  // game event names
  TARGET_HIT = "target:hit"
  TEAM_HIT = "team:hit"

  // game event requests (from game to game engine)
  RANDOM_TEAM_HIT = "random:team:hit" // game requests engine for a random team target hit count

  // node
  NODE_READY = "node:ready"
  NODE_IS_READY = "true"
  NODE_IS_NOT_READY = "false"

  // team names
  BLUE_TEAM = "blue"
  RED_TEAM = "red"
  YELLOW_TEAM = "yellow"
  GREEN_TEAM = "green"
)
