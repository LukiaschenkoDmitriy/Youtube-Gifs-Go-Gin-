package dependency

import (
	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/adapter/http/handler"
	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOAuthHandler(c *config.Config, pool *pgxpool.Pool) *handler.OAuthHandler {
	oauthUseCase := usecase.NewOAuthUC(*c)
	userRepository := repository.NewUserRepository(pool)

	return handler.NewOAuthHandler(oauthUseCase, userRepository)
}

func GetGiphyHandler(c *config.Config) *handler.GiphyHandler {
	giphyUseCase := usecase.NewGiphyUC(c)

	return handler.NewGiphyHandler(giphyUseCase)
}

func GetUserHandler(pool *pgxpool.Pool) *handler.UserHandler {
	userRepository := repository.NewUserRepository(pool)

	return handler.NewUserHandler(userRepository)
}

func GetCommentHandler(pool *pgxpool.Pool) *handler.CommentHandler {
	commentRepository := repository.NewCommentRepository(pool)

	return handler.NewCommentHandler(commentRepository)
}
