
package constants

import (
  "time"
  "errors"
)

var (
  WAIT_TIME = 10 * time.Second
  JOIN_REPLAY = false
  ERR_EXISTING_CONNECTION = errors.New("already connected")
  ERR_NO_AGENT_CONFIG = errors.New("no node agent config")
)
