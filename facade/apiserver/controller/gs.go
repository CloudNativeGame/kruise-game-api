package controller

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/service"
	"github.com/gin-gonic/gin"
)

type GsController struct {
	service.GsService
}

func NewGsController() *GsController {
	return &GsController{
		GsService: *service.GetGsService(),
	}
}

func (g *GsController) GetGameServer(c *gin.Context) {
	namespace, name, ok := getNamespaceNamePathParam(c)
	if !ok {
		return
	}
	gs, err := g.GsService.GetGameServer(namespace, name)
	if err != nil {
		slog.Error("get GameServer failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gs)
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

func (g *GsController) UpdateGameServer(c *gin.Context) {
	namespace, name, ok := getNamespaceNamePathParam(c)
	if !ok {
		return
	}
	jsonPatch, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.Error("get jsonPatch body failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	result, err := g.GsService.UpdateGameServer(namespace, name, jsonPatch)
	if err != nil {
		slog.Error("update GameServer failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	c.JSON(http.StatusOK, result)
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

func (g *GsController) DeleteGameServer(c *gin.Context) {
	namespace, name, ok := getNamespaceNamePathParam(c)
	if !ok {
		return
	}
	result, err := g.GsService.DeleteGameServer(namespace, name)
	if err != nil {
		slog.Error("delete GameServer failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
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
