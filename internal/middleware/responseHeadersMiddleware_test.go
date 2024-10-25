package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derickit/go-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseHeadersMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ResponseHeadersMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, "SAMEORIGIN", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; preload", resp.Header().Get("Strict-Transport-Security"),
		"All expected headers should be set")
}

func TestResponseHeadersMiddleware_CustomHeaders(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ResponseHeadersMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.Writer.Header().Set("Custom-Header", "Custom-Value")
		c.String(http.StatusOK, "Test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, "SAMEORIGIN", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; preload", resp.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "custom-value", resp.Header().Get("Custom-Header"),
		"Custom-Header should be set along with Standard Headers")

}

func TestRespOnseHeadersMiddleware_NoCache(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ResponseHeadersMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.String(http.StatusOK, "Test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the Cache-Control header should not be overwritten
	assert.Equal(t, "no-store", resp.Header().Get("Cache-Control"),
		"Cache-Control should not be overwritten")
}

func TestResponseHeadersMiddleware_NoHeadersSet(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ResponseHeadersMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Empty(t, resp.Header().Get("Custom-Header"), "Custom-Header should not be set")
}
