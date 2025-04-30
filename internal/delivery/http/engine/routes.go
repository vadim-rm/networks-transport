package engine

import (
	"github.com/gin-gonic/gin"

	"transport/internal/delivery/http/handlers"
)

func InitializeExternalRoutes(
	engine *gin.Engine,
	message *handlers.Message,
	segment *handlers.Segment,
) {
	engine.POST("send", message.Send)
	engine.POST("transfer", segment.Receive)
}
