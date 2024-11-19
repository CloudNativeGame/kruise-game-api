package main

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/controller"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/healthz", controller.Healthz)

	gsController := controller.NewGsController()
	r.GET("/v1/gameservers", gsController.GetGameServers)
	r.POST("/v1/gameservers", gsController.UpdateGameServers)
}
