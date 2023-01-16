package baserest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewPingHandler() Handler {
	return &pingHandler{}
}

type pingHandler struct{}

func (pingHandler *pingHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", pingHandler.getPing)
}

func (pingHandler *pingHandler) getPing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
