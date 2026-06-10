package handler

import (
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdditionalHandler struct{}

func NewAdditionalHandler() *AdditionalHandler {
	return &AdditionalHandler{}
}

func (h *AdditionalHandler) Ping(c *gin.Context) {
	c.JSON(200, utils.GetSuccessResponse(nil, "Pong"))
}
