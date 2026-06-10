package middlware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// ANSI colors
const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

func DebugRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		fmt.Printf("%s[REQUEST]%s\n  Method : %s\n  Path   : %s\n  Body   : %s\n",
			colorGreen, colorReset,
			c.Request.Method,
			c.Request.URL.Path,
			string(body),
		)

		c.Next()
	}
}

func DebugResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		fmt.Printf("%s[RESPONSE]%s\n  Status : %d\n  Body   : %s\n",
			colorRed, colorReset,
			c.Writer.Status(),
			blw.body.String(),
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
