package engine

import (
	"github.com/gin-gonic/gin"

	"transport/internal/delivery/http/handlers"
)

func InitializeExternalRoutes(
	engine *gin.Engine,
	message *handlers.Message,
) {
	engine.POST("send", message.Send)
}
