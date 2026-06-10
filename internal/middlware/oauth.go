package middlware

import (
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/errorcode"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

func OAuthUserInMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session.GetUserId(c) != "" {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(errors.New("User already logged in"), errorcode.AlreadyAuthorized))
			c.Abort()
			return
		}

		c.Next()
	}
}

func OAuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session.GetUserId(c) == "" {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(errors.New("you need to log in first"), errorcode.NotAuthorized))
			c.Abort()
			return
		}

		c.Next()
	}
}
