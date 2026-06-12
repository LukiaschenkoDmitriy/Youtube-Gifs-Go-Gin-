package handler

import (
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	oauthUseCase   *usecase.OAuthUC
	userRepository *repository.UserRepository
}

func NewOAuthHandler(oauthUseCase *usecase.OAuthUC, userRepository *repository.UserRepository) *OAuthHandler {
	return &OAuthHandler{
		oauthUseCase:   oauthUseCase,
		userRepository: userRepository,
	}
}

func (oa *OAuthHandler) RedirectToLogin(c *gin.Context) {
	state, err := utils.RandomHex(32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse("Failed to generate random state", err))
		return
	}

	if err := session.SetOAuthState(c, state); err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse("Failed to set OAuth state in session", err))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, oa.oauthUseCase.GetLoginRedirectUrl(state))
}

func (oa *OAuthHandler) Logout(c *gin.Context) {
	if err := session.ClearSession(c); err != nil {
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse("Failed to clear session", err))
		return
	}

	c.JSON(200, utils.GetSuccessResponse(nil, "Logout Success"))
}

func (oa *OAuthHandler) HandleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	cacheStatus, ok := session.GetOAuthState(c)

	if !ok {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "OAuth State not found in session.",
		})
		return
	}

	code, state := c.Query("code"), c.Query("state")

	if cacheStatus != state {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "OAuth State mismatch.",
		})
		return
	}

	if err := session.DeleteOAuthState(c); err != nil {
		c.HTML(http.StatusInternalServerError, "auth-error.html", gin.H{
			"message": "Failed to delete OAuth state from session.",
		})
		return
	}

	user, err := oa.oauthUseCase.GetUser(ctx, code)

	if err != nil {
		message := "Authorization was interrupted or cancelled."

		switch {
		case errors.Is(err, usecase.ErrNoYouTubeChannel):
			message = "A YouTube channel is required to use this application."
		case errors.Is(err, usecase.ErrYouTubeAccessDenied):
			message = "YouTube access was not granted. Please allow it and try again."
		}

		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": message,
		})
		return
	}

	user, err = oa.userRepository.CreateOrGetUser(ctx, user)

	if err != nil {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "Failed to create or get user.",
		})
		return
	}

	if err := session.SetSessionUserId(c, user); err != nil {
		c.HTML(http.StatusInternalServerError, "auth-error.html", gin.H{
			"message": "Failed to save session.",
		})
		return
	}

	c.HTML(http.StatusOK, "auth-success.html", gin.H{})
}
