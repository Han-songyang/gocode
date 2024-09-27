package main

import (
	"github.com/gin-gonic/gin"
	events "webook/internal/event"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
