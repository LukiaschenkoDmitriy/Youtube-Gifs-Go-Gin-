package handler

import (
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
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

	commentId, userId := c.Param("commentId"), c.GetString("userId")

	err := h.cr.DeleteComment(ctx, commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Comment not found", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "delete comment successfully"))
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	ctx := c.Request.Context()

	commentId, userId := c.Param("commentId"), c.GetString("userId")

	comment, err := h.cr.GetCommentById(ctx, commentId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Get comment failed", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comment, "Comment"))
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	ctx := c.Request.Context()

	comment, userId := new(domain.CreateCommentRequest), c.GetString("userId")

	if err := c.ShouldBindJSON(comment); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err.Error(), err))
		return
	}

	nComment, err := h.cr.CreateComment(ctx, comment, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Create comment failed", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nComment, "Comment created"))
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	ctx := c.Request.Context()

	commentId, userId := c.Param("commentId"), c.GetString("userId")

	err := h.cr.LikeComment(ctx, userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Comment not found", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) DislikeComment(c *gin.Context) {
	ctx := c.Request.Context()

	commentId, userId := c.Param("commentId"), c.GetString("userId")

	err := h.cr.DislikeComment(ctx, userId, commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse("Comment not found", err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(nil, "Reaction changed"))
}

func (h *CommentHandler) GetCommentByVideoId(c *gin.Context) {
	ctx := c.Request.Context()

	videoId := c.Param("videoId")

	userId, ok := session.GetUserId(c)

	pagination := new(domain.PaginationRequest)

	if err := c.ShouldBindQuery(pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err.Error(), err))
		return
	}

	if pagination.Cursor < 0 {
		pagination.Cursor = 0
	}

	if pagination.Limit < 10 {
		pagination.Limit = 10
	}

	if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	var userIdPtr *string
	if ok {
		userIdPtr = &userId
	}

	comments, err := h.cr.GetCommentByVideoId(ctx, videoId, userIdPtr, pagination.Cursor, pagination.Limit)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, utils.GetSuccessResponse(comments, "Video comments"))
}
