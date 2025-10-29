package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := InitWebServer()
	app.server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	app.server.Run(":8080")
}
