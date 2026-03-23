package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Ping(c *gin.Context) {
	// use dependency:
	// s.userService.ListUsers()....
	c.JSON(http.StatusOK, "pong")
}
