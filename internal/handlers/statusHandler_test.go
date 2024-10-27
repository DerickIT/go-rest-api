package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derickit/go-rest-api/internal/db/mocks"
	"github.com/derickit/go-rest-api/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func UnMarshalStatusResponse(resp *http.Response) (string, error) {
	body, _ := io.ReadAll(resp.Body)
	var statusResponse string
	err := json.Unmarshal(body, &statusResponse)
	return statusResponse, err
}

func TestStatusSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.PingFunc = func() error {
		return nil
	}
	s := handlers.NewStatusController(&mocks.MockMongoMgr{})

	s.CheckStatus(c)

	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, handlers.UP, statusResponse)

}

func TestStatusDown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.PingFunc = func() error {
		return errors.New("DB connection failed")
	}
	s := handlers.NewStatusController(&mocks.MockMongoMgr{})

	s.CheckStatus(c)

	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}

	assert.EqualValues(t, http.StatusFailedDependency, resp.StatusCode)
	assert.EqualValues(t, handlers.DOWN, statusResponse)
}
