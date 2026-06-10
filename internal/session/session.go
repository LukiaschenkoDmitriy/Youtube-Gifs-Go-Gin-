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

func ClearSession(c *gin.Context) error {
	session := sessions.Default(c)

	session.Clear()
	return session.Save()
}

func GetUserId(c *gin.Context) (string, bool) {
	session := sessions.Default(c)

	userId, ok := session.Get("user_id").(string)

	return userId, ok
}

func SetSessionUserId(c *gin.Context, u *domain.User) error {
	session := sessions.Default(c)

	session.Set("user_id", u.ID)
	return session.Save()
}

func GetOAuthState(c *gin.Context) (string, bool) {
	session := sessions.Default(c)

	oauthState, ok := session.Get("oauth_state").(string)

	return oauthState, ok
}

func SetOAuthState(c *gin.Context, state string) error {
	session := sessions.Default(c)

	session.Set("oauth_state", state)
	return session.Save()
}

func DeleteOAuthState(c *gin.Context) error {
	session := sessions.Default(c)

	session.Delete("oauth_state")
	return session.Save()
}
