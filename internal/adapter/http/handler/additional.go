package handler

import (
	"net/http"

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

func (h *AdditionalHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
