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

type GsController struct {
	filter  *filter.Filter
	updater *updater.Updater
}

func NewGsController() *GsController {
	kubeOption := options.KubeOption{
		KubeConfigPath:          os.Getenv("KUBECONFIG"),
		InformersReSyncInterval: time.Second * 30,
	}
	return &GsController{
		filter: filter.NewFilter(&filter.FilterOption{
			KubeOption: kubeOption,
		}),
		updater: updater.NewUpdater(&updater.UpdaterOption{
			KubeOption: kubeOption,
		}),
	}
}

func (g *GsController) GetGameServers(c *gin.Context) {
	rawFilter := c.Query("filter")
	gameServers, err := g.filter.GetFilteredGameServers(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gameServers)
}

func (g *GsController) UpdateGameServers(c *gin.Context) {
	var request apimodels.UpdateGameServersRequest
	err := c.BindJSON(&request)
	if err != nil {
		slog.Error("update GameServers request body error", "error", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	gameServers, err := g.filter.GetFilteredGameServers(request.Filter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	results := g.updater.Update(gameServers, []byte(request.JsonPatch))
	for _, result := range results {
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, results)
			return
		}
	}

	c.JSON(http.StatusOK, results)
}
