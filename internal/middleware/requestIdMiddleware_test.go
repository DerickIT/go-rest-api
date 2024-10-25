package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derickit/go-rest-api/internal/middleware"
	"github.com/derickit/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReqIDMiddleware(t *testing.T) {
	type reqIDMiddlewareTestCase struct {
		Description  string
		InputReqID   string
		InputReqPath string
	}
	var testCases = []reqIDMiddlewareTestCase{
		{
			Description:  "Test case 1",
			InputReqID:   "123",
			InputReqPath: "/test/1",
		},
		{
			Description:  "Test case 2",
			InputReqID:   "456",
			InputReqPath: "/test/2",
		},
	}
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.ReqIDMiddleware())
	for _, tc := range testCases {
		var hasCorrectReqID bool
		r.GET(tc.InputReqPath, func(ctx *gin.Context) {
			if rID := ctx.Request.Context().Value(util.ContextKey(util.RequestIdentifier)); rID != nil {
				if rIdStr, ok := rID.(string); ok {
					reqIDPassed := len(tc.InputReqID) > 0
					if reqIDPassed && rID == tc.InputReqID || (!reqIDPassed && len(rIdStr) > 0) {
						hasCorrectReqID = true
					}

				}
			}
			ctx.String(200, "OK")
		})

		c.Request, _ = http.NewRequest(http.MethodGet, tc.InputReqPath, nil)
		c.Request.Header.Set(util.RequestIdentifier, tc.InputReqID)
		r.ServeHTTP(resp, c.Request)

		assert.NotEmpty(t, resp.Header().Get(util.RequestIdentifier), tc.Description)

		assert.True(t, hasCorrectReqID, tc.Description)
	}
}
