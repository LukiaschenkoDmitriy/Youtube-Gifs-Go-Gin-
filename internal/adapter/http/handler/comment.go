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
	ctx := c.Request.Context()
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	err := h.cr.DeleteComment(ctx, commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "delete comment successfully"))
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	ctx := c.Request.Context()
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)

	comment, err := h.cr.GetCommentById(ctx, commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comment, "Comment"))
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	comment := new(domain.CreateCommentRequest)

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.WrongJSONDataStructure))
		return
	}

	ctx := c.Request.Context()
	rComment, err := h.cr.CreateComment(ctx, comment, session.GetUserId(c))

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.EntityCreationFailed))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(rComment, "Comment created"))
}

func (h *CommentHandler) GetUserComments(c *gin.Context) {
	userId := session.GetUserId(c)

	ctx := c.Request.Context()
	comments, err := h.cr.GetUserComments(ctx, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comments, "User comments"))
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)
	ctx := c.Request.Context()

	err := h.cr.LikeComment(ctx, userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) DislikeComment(c *gin.Context) {
	commentId := c.Param("commentId")
	userId := session.GetUserId(c)
	ctx := c.Request.Context()

	err := h.cr.DislikeComment(ctx, userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) GetCommentByVideoId(c *gin.Context) {
	videoId := c.Param("videoId")

	userId := session.GetUserId(c)
	ctx := c.Request.Context()
	var userIdPtr *string
	if userId != "" {
		userIdPtr = &userId
	}

	comments, err := h.cr.GetCommentByVideoId(ctx, videoId, userIdPtr)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err, errorcode.DataBaseError))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comments, "Video comments"))
}
