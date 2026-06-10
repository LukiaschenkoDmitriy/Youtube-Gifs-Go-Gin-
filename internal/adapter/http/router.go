package http

import (
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http/handler"
	"github.com/dmytrii/youtube-gifs-chat/internal/middlware"
	"github.com/gin-gonic/gin"
)

func CommentRouters(g *gin.RouterGroup, h *handler.CommentHandler) {
	g.GET("/comments/byVideoId/:videoId", h.GetCommentByVideoId)

	authGroup := g.Group("")
	authGroup.Use(middlware.OAuthRequiredMiddleware())
	{
		{
			authGroup.GET("/comments/:commentId", h.GetComment)
		}
		{
			authGroup.POST("/comments", h.CreateComment)
			authGroup.POST("/comments/:commentId/like", h.LikeComment)
			authGroup.POST("/comments/:commentId/dislike", h.DislikeComment)
		}
		{
			authGroup.DELETE("/comments/:commentId", h.DeleteComment)
		}
	}
}

func OAuthRouters(g *gin.RouterGroup, h *handler.OAuthHandler, authEndpoint string) {
	g.POST(authEndpoint+"/logout", h.Logout)

	guestGroup := g.Group("").Use(middlware.GuestRequredMiddleware())
	{
		guestGroup.GET(authEndpoint+"/login", h.RedirectToLogin)
		guestGroup.GET(authEndpoint+"/callback", h.HandleCallback)
	}
}

func GiphyRouters(g *gin.RouterGroup, h *handler.GiphyHandler) {
	g.Use(middlware.OAuthRequiredMiddleware())
	{
		g.GET("/giphy/trending", h.GetTrendingGifs)
		g.GET("/giphy/search", h.SearchGifs)
	}
}

func UserRouters(g *gin.RouterGroup, h *handler.UserHandler) {
	g.Use(middlware.OAuthRequiredMiddleware())
	{
		g.GET("/users/current", h.GetCurrentUser)
	}
}

func AdditionalRouters(g *gin.RouterGroup, h *handler.AdditionalHandler) {
	g.POST("/u/ping", h.Ping)
}
