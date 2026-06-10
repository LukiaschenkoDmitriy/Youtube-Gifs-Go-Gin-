package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/errorcode"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type GiphyHandler struct {
	giphyUC *usecase.GiphyUC
}

func NewGiphyHandler(giphyUC *usecase.GiphyUC) *GiphyHandler {
	return &GiphyHandler{
		giphyUC: giphyUC,
	}
}

func (g *GiphyHandler) GetTrendingGifs(c *gin.Context) {
	ctx := c.Request.Context()
	gifs, err := g.giphyUC.GetTrendingGifs(ctx, c.DefaultQuery("offset", "0"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(err, errorcode.ClientError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(gifs, "GIFS Trending"))
}

func (g *GiphyHandler) SearchGifs(c *gin.Context) {
	ctx := c.Request.Context()
	gifs, err := g.giphyUC.GetGifsBySearch(ctx, c.Query("search"), c.DefaultQuery("offset", "0"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(err, errorcode.ClientError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(gifs, "GIFS Search"))
}
