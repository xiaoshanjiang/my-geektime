package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/events"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
