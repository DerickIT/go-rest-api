package middleware

import (
	"time"

	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogMiddleware(lgr *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, _ := lgr.WithReqID(c)
		start := time.Now()
		c.Next()
		l.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("path", c.FullPath()).Str("userAgent", c.Request.UserAgent()).
			Int("respStatus", c.Writer.Status()).
			Dur("elapsedMs", time.Since(start)).Send()
	}
}
