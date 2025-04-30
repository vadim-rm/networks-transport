package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"transport/internal/domain"
	"transport/internal/service"
)

type Message struct {
	transport service.Transport
}

func NewMessage(transport service.Transport) *Message {
	return &Message{transport: transport}
}

type message struct {
	Username string    `form:"username"`
	SendTime time.Time `form:"sendTime"`
	Message  string    `form:"message"`
}

func (h *Message) Send(ctx *gin.Context) {
	var request message
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.transport.SendMessage(ctx, domain.Message{
		Username: request.Username,
		SendTime: request.SendTime,
		Message:  request.Message,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
