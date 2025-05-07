package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"io"
	"net/http"
)

type Segment struct {
	writer *kafka.Writer
}

func NewSegment(writer *kafka.Writer) *Segment {
	return &Segment{writer: writer}
}

func (h *Segment) Transfer(ctx *gin.Context) {
	value, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.writer.WriteMessages(ctx, kafka.Message{Value: value})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
