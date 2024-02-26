package game

import (
  "log"
  "fmt"
  "time"
  "strings"
  "context"
  "math/rand"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameEngine struct {
  conf          *config.Config
  gamechan      *GameChannel
  CurrentGame   Game
  CurrentTeams  []string
  CurrentNodes  []string
  *log.Logger
}

func NewGameEngine(cfg *config.Config, gamechan *GameChannel, logger *log.Logger) *GameEngine {
  return &GameEngine{
    conf:         cfg,
    gamechan:     gamechan,
    CurrentGame:  nil,
    CurrentTeams: cfg.Teams,
    CurrentNodes: cfg.Nodes,
    Logger:       log.New(logger.Writer(), "[GAME]: ", logger.Flags()),
  }
}

func (ge *GameEngine) NewGame(mode string) error {
  newgame := NewGame(mode, ge.conf, ge.gamechan, ge.NewScoreboard(), ge.Logger)
  return ge.MountGame(newgame)
}

func (ge *GameEngine) GameOver() bool {
  return ge.CurrentGame == nil || ge.CurrentGame.Completed()
}

func (ge *GameEngine) GameInProgress() bool {
  return ge.CurrentGame != nil && ge.CurrentGame.Running()
}

func (ge *GameEngine) MountGame(g Game) error {
  if ge.GameOver() {
    ge.Printf("loading new game - %s", g)
    ge.CurrentGame = g
  } else {
    return constants.ERR_GAME_RUNNING
  }
  return nil
}

func (ge *GameEngine) StartGame(ctx context.Context) error {
  if ge.GameInProgress() {
    return constants.ERR_GAME_RUNNING
  }

  expect := len(ge.conf.Nodes)

  if err := ge.WaitForNodes(expect, ge.conf.Timeout); err != nil {
    return constants.ERR_NODES_NOT_READY
  }

  ge.Printf("starting game - %s", ge.CurrentGame)
  return ge.CurrentGame.Start(ctx)
}

func (ge *GameEngine) Start(ctx context.Context) error {
  ge.Printf("starting game engine")
  for {
    select {
    case evt := <-ge.gamechan.EventChan:
      ge.Printf("game engine received game event: %s", evt)
    case evt := <-ge.gamechan.RequestChan:
      ge.Printf("game engine received game event request: %s", evt)
      switch evt.Event {
        case constants.RANDOM_TEAM_HIT:
          if err := ge.RandomTeamHit(rand.Intn(5)); err != nil {
            ge.Printf("cannot generate random team hit: %s", err)
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

func (ge *GameEngine) NewScoreboard() map[string]int {
  sb := map[string]int{}

  for _, team := range ge.CurrentTeams {
    sb[team] = 0
  }

  return sb
}

func (ge *GameEngine) WaitForNodes(expect, timeout int) error {
  for {

    // wait for ready
    resp, err := ge.SendQuery(NewGameQuery(constants.NODE_READY, []byte(""), constants.NODE_TAGS))
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

func (ge *GameEngine) SendEvent(e GameEvent) error {
  ge.gamechan.EventChan <- e // tell everyone
//  ge.GameStats.Events = append(ge.GameStats.Events, e.String())
  return nil
}

func (ge *GameEngine) SendQuery(q GameQuery) (map[string][]byte, error) {
  ge.gamechan.QueryChan <- q
  resp := <-q.Response // block for response
  return resp.Answer, resp.Error
}

func (ge *GameEngine) RandomTeamHits() error {
  for i := 1; i <= rand.Intn(10); i++ {
    if err := ge.RandomTeamHit(rand.Intn(5)); err != nil {
      return err
    }
  }
  return nil
}

func (ge *GameEngine) RandomTeamHit(hits int) error {
  node := ge.RandomNode()
  team := ge.RandomTeam()
  evt := strings.Join([]string{node, constants.TEAM_HIT}, constants.SPLIT)
  pay := fmt.Sprintf("%s%s%d", team, constants.SPLIT, hits)
  if err := ge.SendEvent(NewGameEvent(evt, []byte(pay))); err != nil {
    ge.Printf("error sending random team hit %s: %s", team, err)
    return err
  }

  return nil
}

func (ge *GameEngine) RandomTeam() string {
  if len(ge.CurrentTeams) > 0 {
    return ge.CurrentTeams[rand.Intn(len(ge.CurrentTeams))]
  } else {
    return ""
  }
}

func (ge *GameEngine) RandomNode() string {
  if len(ge.CurrentNodes) > 0 {
    return ge.CurrentNodes[rand.Intn(len(ge.CurrentNodes))]
  } else {
    return ""
  }
}

