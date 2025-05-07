package engine

import (
	"github.com/gin-gonic/gin"

	"transport/internal/delivery/http/handlers"
)

func InitializeExternalRoutes(
	engine *gin.Engine,
	segment *handlers.Segment,
) {
	engine.POST("transfer", segment.Transfer)
}
