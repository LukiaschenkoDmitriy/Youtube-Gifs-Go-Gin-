package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/cache"
	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ur    *repository.UserRepository
	Cache *cache.Cache
}

func NewUserHandler(ur *repository.UserRepository, cch *cache.Cache) *UserHandler {
	return &UserHandler{
		ur:    ur,
		Cache: cch,
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

func (uh *UserHandler) UpdateUserSettings(c *gin.Context) {
	ctx := c.Request.Context()

	settings, userId := new(domain.UpdateUserSettingsRequest), c.GetString("userId")

	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err.Error(), err))
		return
	}

	if err := uh.ur.UpdateUserSettings(ctx, userId, *settings); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Failed to update user settings", err))
		return
	}

	uh.Cache.Invalidate(userId + ":/e/users/current")

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "User successfully updated"))

}
