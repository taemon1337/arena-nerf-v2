package game

type GameEvent struct {
  Event         string        `yaml:"event" json:"event"`
  Payload       []byte        `yaml:"payload" json:"payload"`
}

type GameQueryResponse struct {
  Answer          map[string][]byte
  Error           error
}

type GameQuery struct {
  Query           string                  `yaml:"query" json:"query"`
  Payload         []byte                  `yaml:"payload" json:"payload"`
  Tags            map[string]string       `yaml:"tags" json:"tags"`
  Response        chan GameQueryResponse  `yaml:"-" json:"-"`
}

type SensorEvent struct {
  Event         string        `yaml:"event" json:"event"`
  Payload       []byte        `yaml:"payload" json:"payload"`
}

func NewGameEvent(name string, payload []byte) GameEvent {
  return GameEvent{
    Event:    name,
    Payload:  payload,
  }
}

func NewGameQuery(query string, payload []byte, tags map[string]string) GameQuery {
  return GameQuery{
    Query:     query,
    Payload:   payload,
    Tags:      tags,
    Response:  make(chan GameQueryResponse, 0),
  }
}

func NewGameQueryResponse(resp map[string][]byte, err error) GameQueryResponse {
  return GameQueryResponse{
    Answer:   resp,
    Error:    err,
  }
}

func NewSensorEvent(name string, payload []byte) *SensorEvent {
  return &SensorEvent{
    Event:    name,
    Payload:  payload,
  }
}
