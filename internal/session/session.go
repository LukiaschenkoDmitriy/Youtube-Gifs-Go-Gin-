package session

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func CreateStore(c *config.Config) *cookie.Store {
	store := cookie.NewStore([]byte(c.SessionSecret))

	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400 * 7,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return &store
}

func GetUserId(c *gin.Context) string {
	session := sessions.Default(c)

	userId := session.Get("user_id")

	if userId == nil {
		return ""
	}

	return userId.(string)
}

func SetSessionUserId(c *gin.Context, u *domain.User) {
	session := sessions.Default(c)

	session.Set("user_id", u.ID)
	session.Save()
}
