package http

import (
	"time"

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

func OAuthRouters(g *gin.RouterGroup, h *handler.OAuthHandler) {
	g.POST(h.AuthEndpoint+"/logout", h.Logout)

	guestGroup := g.Group("").Use(middlware.GuestRequredMiddleware())
	{
		guestGroup.GET(h.AuthEndpoint+"/login", h.RedirectToLogin)
		guestGroup.GET(h.AuthEndpoint+"/callback", h.HandleCallback)
	}
}

func GiphyRouters(g *gin.RouterGroup, h *handler.GiphyHandler) {
	g.Use(middlware.OAuthRequiredMiddleware())
	{
		tg := g.Group("")
		{
			tg.Use(middlware.CacheMiddleware(h.Cache, time.Duration(30*time.Minute)))
			tg.GET("/giphy/trending", h.GetTrendingGifs)
		}
		g.GET("/giphy/search", h.SearchGifs)
	}
}

func UserRouters(g *gin.RouterGroup, h *handler.UserHandler) {
	g.Use(middlware.OAuthRequiredMiddleware())
	{
		cug := g.Group("")
		{
			cug.Use(middlware.UserCacheMiddleware(h.Cache, time.Duration(time.Hour)))
			cug.GET("/users/current", h.GetCurrentUser)
		}

	}
}

func AdditionalRouters(g *gin.RouterGroup, h *handler.AdditionalHandler) {
	g.POST("/u/ping", h.Ping)
}
