package middlware

import (
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

func GuestRequredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := session.GetUserId(c)

		if ok {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse("User already logged in", errors.New("Authorized")))
			c.Abort()
			return
		}

		c.Next()
	}
}

func OAuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := session.GetUserId(c)

		if !ok {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Unauthorized", errors.New("Unauthorized")))
			c.Abort()
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
