package constants

import (
  "errors"
)

var (
  ERR_GAME_RUNNING = errors.New("current game is still running")
  ERR_NODES_NOT_READY = errors.New("current game nodes are not ready (or not enough are ready")
)
