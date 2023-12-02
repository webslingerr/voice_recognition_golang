package functions

import (
	"project/config"
	"project/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg config.Config
	log logger.LoggerI
}

func NewHandler(cfg config.Config, log logger.LoggerI) *Handler {
	return &Handler{
		cfg: cfg,
		log: log,
	}
}

func (h *Handler) handleResponse(c *gin.Context, status int, data interface{}) {
	switch code := status; {
	case code < 300:
		h.log.Info(
			"response",
			logger.Int("code", status),
		)
	case code < 400:
		h.log.Warn(
			"response",
			logger.Int("code", status),
			logger.Any("data", data),
		)
	default:
		h.log.Error(
			"response",
			logger.Int("code", status),
			logger.Any("data", data),
		)
	}

	var resp interface{}
	if msg, ok := data.(string); ok {
		resp = gin.H{"response": msg}
	} else {
		resp = data
	}

	c.JSON(status, resp)
}
