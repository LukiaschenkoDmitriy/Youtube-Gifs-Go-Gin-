package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ur *repository.UserRepository
}

func NewUserHandler(ur *repository.UserRepository) *UserHandler {
	return &UserHandler{
		ur: ur,
	}
}

func (uh *UserHandler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetString("userId")

	user, err := uh.ur.GetUserById(ctx, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Get user failed", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(user, "User retrieved successfully"))
}
