package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derickit/go-rest-api/internal/db/mocks"
	"github.com/derickit/go-rest-api/internal/handlers"
	"github.com/derickit/go-rest-api/internal/models/data"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewSeedHandler(t *testing.T) {
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{})
	assert.IsType(t, &handlers.SeedHandler{}, sd)
}

func TestSeedDB_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "random-id", nil
		},
	})

	sd.SeedDB(c)
	resp := recorder.Result()
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)

}

func TestSeedDB_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "", assert.AnError
		},
	})

	sd.SeedDB(c)
	resp := recorder.Result()
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)

}
