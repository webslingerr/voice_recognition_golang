package main

import (
	"project/api"
	"project/config"
	"project/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	var loggerLevel = new(string)
	*loggerLevel = logger.LevelDebug
	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)

	}

	log := logger.NewLogger("sas_project", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	r := gin.New()

	// calling logger
	r.Use(gin.Recovery(), gin.Logger())

	api.SetUpApi(r, &cfg, log)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
