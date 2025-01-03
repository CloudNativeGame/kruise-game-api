package controller

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/service"
	"github.com/gin-gonic/gin"
)

type GssController struct {
	service.GssService
}

func NewGssController() *GssController {
	return &GssController{
		GssService: *service.GetGssService(),
	}
}

func (g *GssController) GetGameServerSet(c *gin.Context) {
	namespace, name, ok := getNamespaceNamePathParam(c)
	if !ok {
		return
	}
	gss, err := g.GssService.GetGameServerSet(namespace, name)
	if err != nil {
		slog.Error("get GameServerSet failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gss)
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

func (g *GssController) UpdateGameServerSet(c *gin.Context) {
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
	result, err := g.GssService.UpdateGameServerSet(namespace, name, jsonPatch)
	if err != nil {
		slog.Error("update GameServerSet failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	c.JSON(http.StatusOK, result)
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

func (g *GssController) DeleteGameServerSet(c *gin.Context) {
	namespace, name, ok := getNamespaceNamePathParam(c)
	if !ok {
		return
	}
	result, err := g.GssService.DeleteGameServerSet(namespace, name)
	if err != nil {
		slog.Error("delete GameServerSet failed", "error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	c.JSON(http.StatusOK, result)
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
