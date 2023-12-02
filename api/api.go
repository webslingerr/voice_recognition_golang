package api

import (
	"project/api/docs"
	"project/api/functions"
	"project/config"
	"project/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpApi(r *gin.Engine, cfg *config.Config, logger logger.LoggerI) {
	handler := functions.NewHandler(*cfg, logger)

	docs.SwaggerInfo.Title = cfg.SwaggerTitle
	docs.SwaggerInfo.Version = cfg.SwaggerVersion

	r.POST("/load-recording", handler.LoadRecording)
	r.POST("/recognize-voice", handler.RecognizeVoice)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
