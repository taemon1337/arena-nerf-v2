package constants

var (
  // game modes
  GAME_MODE_SIMULATION = "simulation"

  // game states
  GAME_STATE_INIT = "game:init"
  GAME_STATE_RUNNING = "game:running"
  GAME_STATE_ENDED = "game:over"

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
