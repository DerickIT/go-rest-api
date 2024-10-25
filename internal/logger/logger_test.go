package logger_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models"
	"github.com/derickit/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	env := models.ServiceEnv{Name: "dev"}
	lgr := logger.Setup(env)
	assert.NotNil(t, lgr)
}

func TestWithRegID(t *testing.T) {
	env := models.ServiceEnv{Name: "test"}
	lgr := logger.Setup(env)
	ginCtx := &gin.Context{}
	ginCtx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	_, reqID := lgr.WithReqID(ginCtx)
	assert.Empty(t, reqID)
	reqIDValue := "1234567890"
	ctx := context.WithValue(ginCtx.Request.Context(), util.ContextKey(util.RequestIdentifier), reqIDValue)
	ginCtx.Request = ginCtx.Request.WithContext(ctx)

	_, newReqID := lgr.WithReqID(ginCtx)
	assert.Equal(t, reqIDValue, newReqID)

	ctx = context.WithValue(ginCtx.Request.Context(), util.ContextKey(util.RequestIdentifier), 123)
	ginCtx.Request = ginCtx.Request.WithContext(ctx)
	_, newReqID = lgr.WithReqID(ginCtx)
	assert.Empty(t, newReqID)
}

func TestSetupOnce(t *testing.T) {
	env := models.ServiceEnv{Name: "test"}
	tempFile, err := os.CreateTemp("", "uTest.log")
	require.NoError(t, err)
	defer func(name string) {
		errRemove := os.Remove(name)
		if err != nil {
			t.Log(errRemove)
		}
	}(tempFile.Name())

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			lgr := logger.Setup(env)
			assert.NotNil(t, lgr)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGetZerologLevel(t *testing.T) {
	tests := []struct {
		name       string
		inputLevel string
		expected   zerolog.Level
	}{
		{"debug", "debug", zerolog.DebugLevel},
		{"info", "info", zerolog.InfoLevel},
		{"warn", "warn", zerolog.WarnLevel},
		{"error", "error", zerolog.ErrorLevel},
		{"unknown", "unknown", zerolog.InfoLevel},
		{"fatal", "fatal", zerolog.FatalLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := logger.ZerologLevel(tt.inputLevel)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
