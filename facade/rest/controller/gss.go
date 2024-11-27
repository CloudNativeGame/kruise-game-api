package controller

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type GssController struct {
	filter  *filter.GssFilter
	updater *updater.Updater
}

func NewGssController() *GssController {
	kubeOption := options.KubeOption{
		KubeConfigPath:          os.Getenv("KUBECONFIG"),
		InformersReSyncInterval: time.Second * 30,
	}
	return &GssController{
		filter: filter.NewGssFilter(&filter.FilterOption{
			KubeOption: kubeOption,
		}),
		updater: updater.NewUpdater(&updater.UpdaterOption{
			KubeOption: kubeOption,
		}),
	}
}

func (g *GssController) GetGameServerSets(c *gin.Context) {
	rawFilter := c.Query("filter")
	gameServerSets, err := g.filter.GetFilteredGameServerSets(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gameServerSets)
}

func (g *GssController) UpdateGameServerSets(c *gin.Context) {
	var request apimodels.UpdateGameServerSetsRequest
	err := c.BindJSON(&request)
	if err != nil {
		slog.Error("update GameServerSets request body error", "error", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	gameServerSets, err := g.filter.GetFilteredGameServerSets(request.Filter)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	results := g.updater.UpdateGameServerSets(gameServerSets, []byte(request.JsonPatch))
	for _, result := range results {
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, results)
			return
		}
	}

	c.JSON(http.StatusOK, results)
}
