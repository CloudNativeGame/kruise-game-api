package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Healthz(c *gin.Context) {
	c.String(http.StatusOK, "")
}
