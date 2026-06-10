package main

import (
	config2 "github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http"
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http/handler"
	"github.com/dmytrii/youtube-gifs-chat/internal/dependency"
	"github.com/dmytrii/youtube-gifs-chat/internal/middlware"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	c, r := config2.GetConfig(), gin.Default()
	p, err := usecase.NewDatabase(c)

	r.LoadHTMLGlob("templates/*")

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
	})

	if err != nil {
		panic(err)
	}

	setupGlobalMiddlewares(r, c)
	setupRoutes(r, c, p)

	r.Run()
}

func setupRoutes(e *gin.Engine, cfg *config2.Config, p *pgxpool.Pool) {
	http.OAuthRouters(e.Group(""), dependency.GetOAuthHandler(cfg, p))
	http.GiphyRouters(e.Group(""), dependency.GetGiphyHandler(cfg))

	// Entity Routes
	eg := e.Group("/e")

	http.UserRouters(eg.Group(""), dependency.GetUserHandler(p))
	http.CommentRouters(eg.Group(""), dependency.GetCommentHandler(p))

	http.AdditionalRouters(e.Group(""), handler.NewAdditionalHandler())
}

func setupGlobalMiddlewares(e *gin.Engine, cfg *config2.Config) {
	//e.Use(middlware.JsonAcceptHeaderMiddleware())
	e.Use(middlware.DebugResponseLogger())
	e.Use(middlware.DebugRequestLogger())
	e.Use(sessions.Sessions("session", *session.CreateStore(cfg)))
}
