package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-contrib/sessions"
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
	c.Redirect(http.StatusTemporaryRedirect, oa.oauthUseCase.GetLoginRedirectUrl())
}

func (oa *OAuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.JSON(200, utils.GetSuccessResponse(nil, "Logout Success"))
}

func (oa *OAuthHandler) HandleCallback(c *gin.Context) {
	err := oa.oauthUseCase.InitClient(c.Query("code"))

	if err != nil {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "Authorization was interrupted or cancelled.",
		})
		return
	}

	user, err := oa.oauthUseCase.GetUser()

	if err != nil {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "Failed to get user data.",
		})
		return
	}

	user, err = oa.userRepository.CreateOrGetUser(user)

	if err != nil {
		c.HTML(http.StatusBadRequest, "auth-error.html", gin.H{
			"message": "Failed to create.",
		})
		return
	}

	session.SetSessionUserId(c, user)

	c.HTML(http.StatusOK, "auth-success.html", gin.H{})
}
