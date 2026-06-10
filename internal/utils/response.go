package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetSuccessResponse(data interface{}, message string) *gin.H {
	return &gin.H{
		"message": message,
		"data":    data,
	}
}

func GetErrorResponse(message string, err error) *gin.H {
	fmt.Printf("%s[ERROR RESPONSE]%s\n  Message : %s\n  Error   : %s\n",
		"\033[31m", "\033[0m",
		message,
		err.Error(),
	)

	return &gin.H{
		"error": message,
	}
}
