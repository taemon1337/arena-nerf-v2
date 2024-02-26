package game

type GameChannel struct {
  RequestChan   chan GameEvent  // request chan are events being requested to occur (normally from Game to GameEngine)
  EventChan     chan GameEvent  // event chan are approved live events that have/are occurring
  QueryChan     chan GameQuery  // query chan are sent to all nodes and blocks until all responses received
}

func NewGameChannel() *GameChannel {
  return &GameChannel{
    RequestChan:    make(chan GameEvent, 5),
    EventChan:      make(chan GameEvent, 5),
    QueryChan:      make(chan GameQuery, 5),
  }
}
