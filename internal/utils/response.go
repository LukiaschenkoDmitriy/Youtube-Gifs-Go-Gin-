package utils

import (
	"github.com/dmytrii/youtube-gifs-chat/internal/errorcode"
	"github.com/gin-gonic/gin"
)

func GetSuccessResponse(data interface{}, message string) *gin.H {
	return &gin.H{
		"message": message,
		"data":    data,
	}
}

func GetErrorResponse(err error, errorCode errorcode.ErrorCode) *gin.H {
	return &gin.H{
		"error": err.Error(),
		"code":  errorCode,
	}
}
