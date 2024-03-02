package game

import (
  "log"
  "fmt"
  "time"
  "slices"
  "strings"
  "context"
  "math/rand"
  "encoding/json"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameEngine struct {
  conf                  *config.Config
  gamechan              *GameChannel
  CurrentGame           Game
  CurrentGameState      *GameState
  *log.Logger
}

func NewGameEngine(cfg *config.Config, gamechan *GameChannel, logger *log.Logger) *GameEngine {
  return &GameEngine{
    conf:               cfg,
    gamechan:           gamechan,
    CurrentGame:        nil,
    CurrentGameState:   NewGameState(NewGameConfig(cfg)),
    Logger:             log.New(logger.Writer(), "[GAME]: ", logger.Flags()),
  }
}

func (ge *GameEngine) NewGame(mode string) error {
  newgame := NewGame(mode, ge.conf, ge.gamechan, ge.Logger)
  return ge.MountGame(newgame)
}

func (ge *GameEngine) GameInProgress() bool {
  return ge.CurrentGame != nil && ge.CurrentGameState.Running()
}

func (ge *GameEngine) MountGame(g Game) error {
  if ge.GameInProgress() {
    return constants.ERR_GAME_RUNNING
  }
  ge.Printf("loading new game - %s", g)
  // TODO: save old game
  ge.CurrentGame = g
  ge.CurrentGameState = NewGameState(NewGameConfig(ge.conf))
  return nil
}

func (ge *GameEngine) StartGame(ctx context.Context) error {
  if ge.GameInProgress() {
    return constants.ERR_GAME_RUNNING
  }

  expect := len(ge.conf.Nodes)

  if err := ge.WaitForNodes(expect, ge.conf.Timeout); err != nil {
    ge.Printf("error waiting for game nodes: %s", err)
    return constants.ERR_NODES_NOT_READY
  }

  if err := ge.WaitForGameModeSetup(ge.CurrentGame.Mode()); err != nil {
    ge.Printf("error sending game mode to nodes: %s", err)
    return err
  }

  if err := ge.CurrentGameState.GameSetup(); err != nil {
    ge.Printf("error setting game up: %s", err)
    return err
  }

  // validate the number of teams and nodes (fail game instead of return error)
  if err := ge.CurrentGameState.ValidateNodes(); err != nil {
    ge.Printf("error validating node/team configuration: %s", err)
    if err := ge.FailGame(err); err != nil {
      return err
    }
    return nil // returning error will shutdown which we don't want
  }

  ge.Printf("%s:\n---\n%s", ge.CurrentGame, ge.CurrentGameState)

  ge.Printf("passing control to game - %s", ge.CurrentGame)
  return ge.CurrentGame.Start(ctx)
}

func (ge *GameEngine) Start(ctx context.Context) error {
  ge.Printf("starting game engine")

  for {
    // we want to eval the game each loop, not just when an event occurs
    if ge.GameInProgress() {
      ge.Printf(ge.CurrentGameState.GameStatus())
      if ge.CurrentGameState.TimeExpired() {
        ge.Printf("game time expired, ending game")
        if err := ge.EndGame(); err != nil {
          return err
        }
      }

      if time.Since(ge.CurrentGameState.Lastcheck) > (10 * time.Second) {
        if !ge.CurrentGameState.Checking() {
          ge.Printf("checking on scores")
          ge.CurrentGameState.SetChecking(true)

          scoreboard, nodeboard, err := ge.GetScoreboard()
          if err != nil {
            ge.Printf("error compiling node scores: %s", err)
            return err
          }

          ge.CurrentGameState.SetBoards(scoreboard, nodeboard)
        }
      }

      if ge.CurrentGameState.WinningScoreReached() {
        ge.Printf("the winning score has been reached, ending game")
        if err := ge.EndGame(); err != nil {
          return err
        }
      }
    }

    select {
    // game engine only listens to RequestChan (controller listens to EventChan)
    // as it is a shared channel, they would compete over it
    case evt := <-ge.gamechan.RequestChan:
      ge.Printf("game engine received game event request: %s", evt)
      switch evt.Event {
        case constants.GAME_ACTION_BEGIN:
          ge.Printf("game engine requested start game - %s", string(evt.Payload))
          if err := ge.SendEventToNodes(evt); err != nil {
            ge.Printf("error telling nodes to start game: %s", err)
            return err
          }
        case constants.GAME_ACTION_END:
          ge.Printf("game engine requested end game - %s", string(evt.Payload))
          if err := ge.EndGame(); err != nil {
            return err
          }
        case constants.RANDOM_TEAM_HIT:
          if ge.GameInProgress() {
            if err := ge.RandomTeamHit(rand.Intn(5)+1); err != nil {
              ge.Printf("cannot generate random team hit: %s", err)
            }
          } else {
            ge.Printf("game engine received request when no game in progress")
          }
        case constants.RANDOM_SENSOR_COLOR:
          if ge.GameInProgress() {
            if err := ge.RandomSensorColor(); err != nil {
              ge.Printf("cannot generate random team color: %s", err)
            }
          } else {
            ge.Printf("game engine received request when no game in progress")
          }
        case constants.RANDOM_SENSOR_HIT:
          if ge.GameInProgress() {
            if err := ge.RandomSensorHit(1); err != nil {
              ge.Printf("cannot generate random sensor hit: %s", err)
            }
          } else {
            ge.Printf("game engine received request when no game in progress")
          }
        default:
          ge.Printf("Unsupported game event request: %s", evt.Event)
      }
    case <-ctx.Done():
      ge.Printf("stopping game engine")
      return ctx.Err()
    default:
      time.Sleep(3 * time.Second) // do something later
    }
  }
}

func (ge *GameEngine) WaitForNodes(expect, timeout int) error {
  ge.Printf("waiting for nodes to be ready")
  for {
    // wait for ready
    resp, err := ge.SendQueryToNodes(NewGameQuery(constants.NODE_READY, []byte(""), constants.NODE_TAGS))
    if err != nil {
      ge.Printf("error query readiness of nodes: %s", err)
      return err
    }
    
    readycount := 0
    
    for _, val := range resp {
      if string(val) == constants.NODE_IS_READY {
        readycount += 1
      }
    }
    
    if readycount >= expect {
      ge.Printf("nodes ready: %d", readycount)
      break // got expected amount node responses indicating readiness
    } else {
      ge.Printf("waiting for %d ready nodes [%d/%d]...", expect, readycount, expect)
      time.Sleep(time.Duration(timeout))
    }
  }
  return nil
}

func (ge *GameEngine) WaitForGameModeSetup(mode string) error {
  ge.Printf("waiting for all game nodes to have proper configuration")

  // set game mode
  ge.Printf("setting game mode")
  if err := ge.SendEventToNodes(NewGameEvent(constants.GAME_MODE, []byte(mode))); err != nil {
    return err
  }

  // query all nodes game mode
  resp, err := ge.SendQueryToNodes(NewGameQuery(constants.GAME_MODE, []byte(""), constants.NODE_TAGS))
  if err != nil {
    ge.Printf("error querying game node: %s", err)
    return err
  }

  passed := 0

  // check the game mode on each node was properly set
  for node, val := range resp {
    ge.CurrentGameState.AddNode(node)
    if string(val) == mode {
      passed += 1
    } else {
      ge.Printf("node %s game mode was set to %s, not %s", node, val, mode)
    }
  }

  if passed == 0 || passed != len(resp) {
    return ge.WaitForGameModeSetup(mode) // retry until successful
  }

  // send all teams to nodes
  if err := ge.SendEventToNodes(NewGameEvent(constants.GAME_TEAMS, []byte(ge.CurrentGameState.TeamList()))); err != nil {
    return err
  }

  return nil
}

func (ge *GameEngine) FailGame(err error) error {
  ge.CurrentGameState.SetStatus(constants.GAME_STATUS_FAILED)
  ge.CurrentGameState.LogGameEvent(NewGameEvent(constants.GAME_ERROR, []byte(fmt.Sprintf("%s", err))))

  if err := ge.LogGame(); err != nil {
    return err
  }

  ge.CurrentGame = nil
  ge.CurrentGameState = NewGameState(NewGameConfig(ge.conf))
  return nil
}

func (ge *GameEngine) EndGame() error {
  ge.CurrentGameState.EndGame()

  if err := ge.SendEventToGame(NewGameEvent(constants.GAME_ACTION_OFF, []byte("turn off game"))); err != nil {
    ge.Printf("error shutting down game: %s", err)
    // don't need to return here
  }

  if err := ge.SendEventToNodes(NewGameEvent(constants.GAME_ACTION_END, []byte("The game has ended."))); err != nil {
    ge.Printf("error sending game ended event: %s", err)
    return err
  }

  scoreboard, nodeboard, err := ge.GetScoreboard()
  if err != nil {
    ge.Printf("error compiling node scores: %s", err)
    return err
  }

  ge.CurrentGameState.SetBoards(scoreboard, nodeboard)
  ge.Printf("Final Score: %s", scoreboard)

  for team, count := range scoreboard {
    if count > ge.CurrentGameState.Highscore {
      ge.CurrentGameState.SetWinner(team, count)
    }
  }

  ge.Printf("The winning team is %s with a score of %d", ge.CurrentGameState.Winner, ge.CurrentGameState.Highscore)

  if err := ge.SendEventToNodes(NewGameEvent(constants.GAME_WINNER, []byte(ge.CurrentGameState.Winner))); err != nil {
    log.Printf("error sending team winner: %s", err)
    return err
  }

  if err := ge.LogGame(); err != nil {
    return err
  }

  return nil
}

func (ge *GameEngine) LogGame() error {
  if ge.conf.Logdir != "" {
    err := ge.Logstats()
    if err != nil {
      ge.Printf("could not write game log: %s", err)
      return err
    } else {
      ge.Printf("saved game log to %s", ge.Logfile())
    }
  }
  return nil
}

func (ge *GameEngine) GetScoreboard() (map[string]int, map[string]int, error) {
  scoreboard := map[string]int{}
  nodeboard := map[string]int{}
  nodes := ge.CurrentGameState.Nodes
  teams := ge.CurrentGameState.Teams

  resp, err := ge.SendQueryToNodes(NewGameQuery(constants.NODE_SCOREBOARD, []byte(""), constants.NODE_TAGS))
  if err != nil {
    ge.Printf("error querying node scoreboards: %s", err)
    return scoreboard, nodeboard, err
  }

  ge.Printf("SCOREBOARDS: %s", resp)

  // accumulate each node response
  for node, val := range resp {
    nodehits := map[string]int{}

    if err := json.Unmarshal(val, &nodehits); err != nil {
      ge.Printf("cannot parse node hits: %s", err)
    } else {
      for key,count := range nodehits {
        isnode := slices.Contains(nodes, key)
        isteam := slices.Contains(teams, key)

        if isteam {
          if _, ok := scoreboard[key]; ok {
            scoreboard[key] += count
          } else {
            scoreboard[key] = count
          }
        }

        if isnode {
          if _, ok := nodeboard[key]; ok {
            nodeboard[key] += count
          } else {
            nodeboard[key] = count
          }
        }

        if !isteam && !isnode {
          msg := fmt.Sprintf("unrecognized team|node %s found in response from node %s", key, node)
          ge.Printf(msg)
          //ge.GameStats.Events = append(ge.GameStats.Events, msg)
        }
      }
    }
  }

  return scoreboard, nodeboard, nil
}


func (ge *GameEngine) SendEventToGame(e GameEvent) error {
  ge.gamechan.GameChan <- e
//  ge.GameStats.Events = append(ge.GameStats.Events, e.String())
  return nil
}

func (ge *GameEngine) SendEventToNodes(e GameEvent) error {
  ge.gamechan.NodeChan <- e
  return nil
}

func (ge *GameEngine) SendQueryToNodes(q GameQuery) (map[string][]byte, error) {
  ge.gamechan.QueryChan <- q
  resp := <-q.Response // block for response
  return resp.Answer, resp.Error
}

func (ge *GameEngine) RandomTeamHits() error {
  for i := 1; i <= rand.Intn(10); i++ {
    if err := ge.RandomTeamHit(rand.Intn(5) + 1); err != nil {
      return err
    }
  }
  return nil
}

func (ge *GameEngine) RandomTeamHit(hits int) error {
  node := ge.CurrentGameState.RandomNode()
  team := ge.CurrentGameState.RandomTeam()
  evt := strings.Join([]string{node, constants.TEAM_HIT}, constants.SPLIT)
  pay := fmt.Sprintf("%s%s%d", team, constants.SPLIT, hits)
  if err := ge.SendEventToNodes(NewGameEvent(evt, []byte(pay))); err != nil {
    ge.Printf("error sending random team hit %s: %s", team, err)
    return err
  }

  return nil
}

func (ge *GameEngine) RandomSensorHit(hits int) error {
  node := ge.CurrentGameState.RandomNode()
  sensorid := constants.RANDOM_SENSOR_ID
  evt := strings.Join([]string{node, constants.SENSOR_HIT_REQUEST}, constants.SPLIT)
  pay := fmt.Sprintf("%s%s%d", sensorid, constants.SPLIT, hits)
  if err := ge.SendEventToNodes(NewGameEvent(evt, []byte(pay))); err != nil {
    ge.Printf("error sending random sensor hit %s: %s", sensorid, err)
    return err
  }

  return nil
}

func (ge *GameEngine) RandomSensorColor() error {
  node := ge.CurrentGameState.RandomNode()
  sensorid := constants.RANDOM_SENSOR_ID
  evt := strings.Join([]string{node, constants.SENSOR_COLOR_REQUEST}, constants.SPLIT)
  pay := strings.Join([]string{sensorid, constants.RANDOM_COLOR_ID}, constants.SPLIT)
  if err := ge.SendEventToNodes(NewGameEvent(evt, []byte(pay))); err != nil {
    ge.Printf("error sending random sensor color %s: %s", sensorid, err)
    return err
  }

  return nil
}
