package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/cache"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type GiphyHandler struct {
	giphyUC *usecase.GiphyUC
	Cache   *cache.Cache
}

func NewGiphyHandler(giphyUC *usecase.GiphyUC, cch *cache.Cache) *GiphyHandler {
	return &GiphyHandler{
		giphyUC: giphyUC,
		Cache:   cch,
	}
}

func (g *GiphyHandler) GetTrendingGifs(c *gin.Context) {
	ctx := c.Request.Context()

	gifs, err := g.giphyUC.GetTrendingGifs(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse("Failed to get trending gifs", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(gifs, "GIFS Trending"))
}

func (g *GiphyHandler) SearchGifs(c *gin.Context) {
	ctx := c.Request.Context()

	gifs, err := g.giphyUC.GetGifsBySearch(ctx, c.Query("search"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse("Failed to search gifs", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(gifs, "GIFS Search"))
}
