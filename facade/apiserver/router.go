package main

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/healthz", controller.Healthz)

	gsController := controller.NewGsController()
	r.GET("/v1/gameservers", gsController.GetGameServers)
	r.POST("/v1/gameservers", gsController.UpdateGameServers)
	r.DELETE("/v1/gameservers", gsController.DeleteGameServers)

	gssController := controller.NewGssController()
	r.GET("/v1/gameserversets", gssController.GetGameServerSets)
	r.POST("/v1/gameserversets", gssController.UpdateGameServerSets)
	r.DELETE("/v1/gameserversets", gssController.DeleteGameServerSets)
}
