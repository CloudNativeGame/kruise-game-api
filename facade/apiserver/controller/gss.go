package controller

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type GssController struct {
	service.GssService
}

func NewGssController() *GssController {
	return &GssController{
		GssService: *service.GetGssService(),
	}
}

func (g *GssController) GetGameServerSets(c *gin.Context) {
	rawFilter := c.Query("filter")
	gameServerSets, err := g.GssService.GetGameServerSets(rawFilter)
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

	results, err := g.GssService.UpdateGameServerSets(&request)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
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

func (g *GssController) DeleteGameServerSets(c *gin.Context) {
	rawFilter := c.Query("filter")
	results, err := g.GssService.DeleteGameServerSets(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
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
