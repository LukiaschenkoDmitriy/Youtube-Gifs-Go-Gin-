package main

import (
	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http"
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http/handler"
	"github.com/dmytrii/youtube-gifs-chat/internal/database"
	"github.com/dmytrii/youtube-gifs-chat/internal/dependency"
	"github.com/dmytrii/youtube-gifs-chat/internal/middlware"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	c, r, p := initDeps()

	r.LoadHTMLGlob("templates/*")

	setupErrorRoute(r)
	setupGlobalMiddlewares(r, c)
	setupRoutes(r, c, p)

	r.Run()
}

func setupErrorRoute(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
	})
}

func setupRoutes(e *gin.Engine, c *config.Config, p *pgxpool.Pool) {
	http.OAuthRouters(e.Group(""), dependency.GetOAuthHandler(c, p), c.AuthEndpoint)
	http.GiphyRouters(e.Group(""), dependency.GetGiphyHandler(c))

	// Entity Routes
	eg := e.Group("/e")

	http.UserRouters(eg.Group(""), dependency.GetUserHandler(p))
	http.CommentRouters(eg.Group(""), dependency.GetCommentHandler(p))

	http.AdditionalRouters(e.Group(""), handler.NewAdditionalHandler())
}

func setupGlobalMiddlewares(e *gin.Engine, c *config.Config) {
	//e.Use(middlware.JsonAcceptHeaderMiddleware())
	e.Use(middlware.DebugResponseLogger())
	e.Use(middlware.DebugRequestLogger())
	e.Use(sessions.Sessions("session", *session.CreateStore(c)))
}

func initDeps() (*config.Config, *gin.Engine, *pgxpool.Pool) {
	c, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	p, err := database.NewDatabase(c)

	if err != nil {
		panic(err)
	}

	return c, gin.Default(), p
}
