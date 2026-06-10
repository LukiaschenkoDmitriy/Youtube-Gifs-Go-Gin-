package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/dmytrii/youtube-gifs-chat/internal/errorcode"
	"github.com/dmytrii/youtube-gifs-chat/internal/repository"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	cr *repository.CommentRepository
}

func NewCommentHandler(cr *repository.CommentRepository) *CommentHandler {
	return &CommentHandler{
		cr: cr,
	}
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	err := h.cr.DeleteComment(commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "delete comment successfully"))
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	comment, err := h.cr.GetCommentById(commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comment, "Comment"))
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	comment := new(domain.Comment)

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.WrongJSONDataStructure))
		return
	}

	comment, err := h.cr.CreateComment(comment)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.EntityCreationFailed))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(&comment, "Comment created"))
}

func (h *CommentHandler) GetUserComments(c *gin.Context) {
	userId := session.GetUserId(c)

	comments, err := h.cr.GetUserComments(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comments, "User comments"))
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	err := h.cr.LikeComment(userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) DislikeComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	err := h.cr.DislikeComment(userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) GetCommentByVideoId(c *gin.Context) {
	videoId := c.Param("videoId")

	userId := session.GetUserId(c)

	var userIdPtr *string
	if userId != "" {
		userIdPtr = &userId
	}

	comments, err := h.cr.GetCommentByVideoId(videoId, userIdPtr)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comments, "Video comments"))
}
