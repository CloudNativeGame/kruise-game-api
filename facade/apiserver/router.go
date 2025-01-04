package main

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/controller"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/healthz", controller.Healthz)

	gsController := controller.NewGsController()
	r.GET("/v1/gameservers", gsController.GetGameServers)
	r.GET("/v1/gameserver/:namespace/:name", gsController.GetGameServer)

	r.POST("/v1/gameservers", gsController.UpdateGameServers)
	r.PATCH("/v1/gameserver/:namespace/:name", gsController.UpdateGameServer)

	r.DELETE("/v1/gameservers", gsController.DeleteGameServers)
	r.DELETE("/v1/gameserver/:namespace/:name", gsController.DeleteGameServer)

	gssController := controller.NewGssController()
	r.GET("/v1/gameserversets", gssController.GetGameServerSets)
	r.GET("/v1/gameserverset/:namespace/:name", gssController.GetGameServerSet)

	r.POST("/v1/gameserversets", gssController.UpdateGameServerSets)
	r.PATCH("/v1/gameserverset/:namespace/:name", gssController.UpdateGameServerSet)

	r.DELETE("/v1/gameserversets", gssController.DeleteGameServerSets)
	r.DELETE("/v1/gameserverset/:namespace/:name", gssController.DeleteGameServerSet)
}
