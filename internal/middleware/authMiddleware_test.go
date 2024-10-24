package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derickit/go-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "status code should be 200")

}

func TestAuthMiddleware_WithNext(t *testing.T) {
	router := gin.New()
	var nextCalled bool
	router.Use(middleware.AuthMiddleware())

	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		c.String(http.StatusOK, "Test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	assert.True(t, nextCalled, "next should be called")
}
