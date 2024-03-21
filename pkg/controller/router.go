package controller

import (
  "os"
  "log"
  "fmt"
  "strings"
  "net/http"
  "path/filepath"
  "github.com/gin-gonic/gin"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

type PayloadForm struct {
  Payload       string      `yaml:"payload" json:"payload"`
}

func (ctrl *Controller) Router() {
  api := ctrl.server.Router.Group("api")
  v1 := api.Group("v1")
  {
    v1.GET("/games/:uuid", ctrl.ApiGameStats())
    v1.POST("/do/:action", ctrl.ApiAction())
  }
}

func (ctrl *Controller) ApiGameStats() func (*gin.Context) {
  return func (c *gin.Context) {
    uuid := c.Param("uuid")
    switch uuid {
      case "all":
        // send listing of all archived game logs
        entries, err := os.ReadDir(ctrl.conf.Logdir)
        if err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": err})
          return
        }
        uuids := []string{"current"}
        for _, e := range entries {
          if strings.HasSuffix(e.Name(), ".json") {
            uuids = append(uuids, strings.TrimSuffix(e.Name(), ".json"))
          }
        }
        c.JSON(http.StatusOK, gin.H{
          "games": uuids,
        })
        return
      case "current":
        if !ctrl.engine.GameInProgress() {
          c.JSON(http.StatusOK, gin.H{
            "msg": "no active game",
          })
        } else {
          // send current game stats
          c.JSON(http.StatusOK, gin.H{
            "stats": ctrl.engine.CurrentGameState,
          })
        }
        return
      default:
        // send archived log file
        http.ServeFile(c.Writer, c.Request, filepath.Join(ctrl.conf.Logdir, uuid + ".json"))
        return
    }
  }
}

func (ctrl *Controller) ApiAction() func (*gin.Context) {
  return func (c *gin.Context) {
    action := c.Param("action")

    if !ctrl.conf.EnableApiActions {
      c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", constants.ERR_API_ACTIONS_NOT_ALLOWED)})
      return
    }

    var payForm PayloadForm
    c.ShouldBindJSON(&payForm)
    err := ctrl.ActionFromUi(action, payForm.Payload)

    if err != nil {
      log.Printf("cannot perform api action %s: %s", action, err)
      c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", err)})
      return
    }

    c.JSON(http.StatusOK, gin.H{
      "message": fmt.Sprintf("%s action sent", action),
    })
  }
}

func (ctrl *Controller) ActionFromUi(action, payload string) error {
  switch action {
    case "ui:game:mode":
      if ctrl.engine.GameInProgress() {
        return constants.ERR_ONGOING_GAME
      } 

      return ctrl.engine.SendEventToNodes(game.NewGameEvent(constants.GAME_MODE, []byte(payload)))
// TODO: how to properly async start the game and attach to waitgroup
//    case "ui:game:begin":
//      go ctrl.game.Run(ctrl.conf.ExpectNodes, ctrl.conf.Timeout)
//      return ctrl.game.SendAction(constants.GAME_ACTION_BEGIN, "web: Start the game!")
    case "ui:game:end":
      return ctrl.engine.EndGame()
    default:
      return constants.ERR_UI_ACTION_NOT_ALLOWED
  }
}
