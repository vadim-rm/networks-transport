package engine

import "github.com/gin-gonic/gin"

func Initialize() *gin.Engine {
	router := gin.Default()

	return router
}
