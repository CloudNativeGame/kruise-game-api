package controller

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type GsController struct {
	service.GsService
}

func NewGsController() *GsController {
	return &GsController{
		GsService: *service.GetGsService(),
	}
}

func (g *GsController) GetGameServers(c *gin.Context) {
	rawFilter := c.Query("filter")
	gameServers, err := g.GsService.GetGameServers(rawFilter)
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

	results, err := g.GsService.UpdateGameServers(&request)
	if err != nil {
		slog.Error("update GameServers failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, result := range results {
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, results)
			return
		}
	}

	c.JSON(http.StatusOK, results)
}

func (g *GsController) DeleteGameServers(c *gin.Context) {
	rawFilter := c.Query("filter")

	results, err := g.GsService.DeleteGameServers(rawFilter)
	if err != nil {
		slog.Error("delete GameServers failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, result := range results {
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, results)
			return
		}
	}

	c.JSON(http.StatusOK, results)
}
