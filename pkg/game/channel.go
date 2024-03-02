package game

type GameChannel struct {
  RequestChan   chan GameEvent  // request chan are events being requested to occur (normally from Game to GameEngine)
  GameChan      chan GameEvent  // game chan are events being sent from GameEngine to Game
  NodeChan      chan GameEvent  // node chan are events being sent from GameEngine to Nodes (over network)
  QueryChan     chan GameQuery  // query chan are sent to all nodes and blocks until all responses received
}

func NewGameChannel() *GameChannel {
  return &GameChannel{
    RequestChan:    make(chan GameEvent, 5),
    GameChan:       make(chan GameEvent, 5),
    NodeChan:       make(chan GameEvent, 5),
    QueryChan:      make(chan GameQuery, 5),
  }
}
