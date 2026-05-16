package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "Driving Authority API",
		"version": "1.0",
		"docs":    "/swagger/index.html",
		"health":  "/api/v1/health",
		"api":     "/api/v1",
	})
}
