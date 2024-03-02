package constants

import (
  "errors"
)

var (
  ERR_GAME_RUNNING = errors.New("current game is still running")
  ERR_NODES_NOT_READY = errors.New("current game nodes are not ready (or not enough are ready")
  ERR_MIN_NODE_COUNT = errors.New("not enough nodes")
  ERR_MAX_NODE_COUNT = errors.New("too many nodes")
  ERR_MIN_TEAM_COUNT = errors.New("not enough teams")
  ERR_MAX_TEAM_COUNT = errors.New("too many teams")
  ERR_REQUIRED_NODE = errors.New("missing required node")
  ERR_REQUIRED_TEAM = errors.New("missing required team")
  ERR_INVALID_NODE_HIT = errors.New("invalid node hit payload - must be <sensor-id>:<sensor-color>:<hit-count>")
)
