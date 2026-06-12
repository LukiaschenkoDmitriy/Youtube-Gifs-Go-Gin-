package middlware

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/dmytrii/youtube-gifs-chat/internal/cache"
	"github.com/dmytrii/youtube-gifs-chat/internal/session"
	"github.com/gin-gonic/gin"
)

type cachedResponse struct {
	Status      int
	ContentType string
	Body        []byte
}

func (r cachedResponse) String() string {
	body := string(r.Body)
	if len(body) > 300 {
		body = body[:300] + "… (truncated)"
	}
	return fmt.Sprintf("{Status: %d, ContentType: %s, Body: %s}", r.Status, r.ContentType, body)
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheMiddleware caches responses by request URI, shared between all users.
func CacheMiddleware(cch *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return cacheMiddleware(cch, ttl, func(c *gin.Context) (string, bool) {
		return c.Request.URL.RequestURI(), true
	})
}

// UserCacheMiddleware caches responses per user, prefixing the key with the
// session userId. Requests without a user in the session are not cached.
func UserCacheMiddleware(cch *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return cacheMiddleware(cch, ttl, func(c *gin.Context) (string, bool) {
		userId, ok := session.GetUserId(c)
		if !ok {
			return "", false
		}
		return userId + ":" + c.Request.URL.RequestURI(), true
	})
}

func cacheMiddleware(cch *cache.Cache, ttl time.Duration, keyFunc func(*gin.Context) (string, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		key, ok := keyFunc(c)
		if !ok {
			c.Next()
			return
		}

		if val, found := cch.Get(key); found {
			resp := val.(cachedResponse)
			c.Data(resp.Status, resp.ContentType, resp.Body)
			c.Abort()
			return
		}

		bw := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = bw

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			cch.Set(key, cachedResponse{
				Status:      c.Writer.Status(),
				ContentType: c.Writer.Header().Get("Content-Type"),
				Body:        bw.body.Bytes(),
			}, ttl)
		}
	}
}
