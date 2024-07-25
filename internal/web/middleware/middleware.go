package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ernestngugi/medvice-backend/internal/web/contexthelper"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIdHeaderKey = "X-Request-ID"
	userAgentHeaderKey = "user-agent"
)

func DefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{

		secureMiddleware(),
		compressMiddleware(),
		corsMiddleware(),

		setRequestIdMiddleware(),
		setupContextMiddleware(),

		panicRecoverMiddleware(),
	}
}

func setupContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		userAgent := c.Request.Header.Get(userAgentHeaderKey)
		ctx = contexthelper.WithUserAgent(ctx, userAgent)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func compressMiddleware() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

func secureMiddleware() gin.HandlerFunc {
	return secure.Secure(secure.Options{
		SSLRedirect:          strings.ToLower(os.Getenv("FORCE_SSL")) == "true",
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	})
}

func setRequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		ctx := contexthelper.WithRequestId(c.Request.Context(), requestId)
		c.Request = c.Request.WithContext(ctx)
		c.Header(RequestIdHeaderKey, requestId)
		c.Next()
	}
}

func panicRecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("recover from panic err %v", err)
				debug.PrintStack()
				fmt.Fprintf(
					c.Writer,
					`{"error_message":"internal server error (%s)"}`,
					contexthelper.RequestId(c.Request.Context()),
				)
			}
		}()
		c.Next()
	}
}
