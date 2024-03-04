package common

import (
  "fmt"
  "errors"
  "strings"
  "strconv"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

func ParsePayload(payload []byte) []string {
  return strings.Split(string(payload), constants.SPLIT)
}

func ParseInt(s string) (int, error) {
  return strconv.Atoi(s)
}

func ParseHit(payload []byte) (string, int, error) {
  parts := ParsePayload(payload)
  if len(parts) < 2 {
    return "", 0, errors.New(fmt.Sprintf("cannot parse hit from %s - should be <key>:<hit-count>", string(payload)))
  }

  hits, err := ParseInt(parts[1])
  if err != nil {
    return "", 0, errors.New(fmt.Sprintf("cannot parse hit from %s - %s", string(payload), err))
  }

  return parts[0], hits, nil
}

func ParseSensorHit(payload []byte) (string, int, error) {
  return ParseHit(payload)
}

func ParseTeamHit(payload []byte) (string, int, error) {
  return ParseHit(payload)
}

func ParseNodeHitPayload(payload []byte) (string, string, int, error) {
  parts := ParsePayload(payload)

  // <sensor-id>:<sensor-color>:<hit-count>
  if len(parts) != 3 {
    return "", "", 0, errors.New(fmt.Sprintf("cannot parse node hit from %s - should be <sensor-id>:<sensor-color>:<hit-count>", string(payload)))
  }

  sensorid := parts[0]
  sensorcolor := parts[1]
  hitcount, err := ParseInt(parts[2])
  if err != nil {
    return "", "", 0, errors.New(fmt.Sprintf("cannot parse node hit count from %s - should be <sensor-id>:<sensor-color>:<hit-count>", string(payload)))
  }

  return sensorid, sensorcolor, hitcount, nil
}

