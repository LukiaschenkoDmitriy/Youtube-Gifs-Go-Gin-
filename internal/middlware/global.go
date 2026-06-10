package middlware

import (
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

func JsonAcceptHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Accept") != "application/json" {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse("json accept header not json", errors.New("Invalid Accept header")))
			return
		}

		c.Next()
	}
}
