package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym0818/mywe/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
