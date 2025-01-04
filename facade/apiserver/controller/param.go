package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getNamespaceNamePathParam(c *gin.Context) (namespace, name string, ok bool) {
	namespace = c.Param("namespace")
	name = c.Param("name")
	if namespace == "" || name == "" {
		msg := "namespace and name must been provided in path parameter"
		slog.Error(msg)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	ok = true
	return
}
