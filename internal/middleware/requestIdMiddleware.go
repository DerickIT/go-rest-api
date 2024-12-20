package middleware

import (
	"context"

	"github.com/derickit/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ReqIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(util.RequestIdentifier)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), util.ContextKey(util.RequestIdentifier), reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(util.RequestIdentifier, reqID)
		c.Next()
	}
}
