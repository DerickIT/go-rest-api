package server_test

import (
	"net/http"
	"testing"

	"github.com/derickit/go-rest-api/internal/db/mocks"
	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models"
	"github.com/derickit/go-rest-api/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListOfRoutes(t *testing.T) {
	svcInfo := models.ServiceEnv{
		Name: "test",
		Port: "8080",
	}
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	router := server.WebRouter(svcInfo, &mocks.MockMongoMgr{}, lgr)
	list := router.Routes()
	mode := gin.Mode()
	assert.Equal(t, gin.ReleaseMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/status",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/seedDB",
	})

	assertRouteNotPresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/ecommerce/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/ecommerce/v1/orders/:id",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/ecommerce/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodDelete,
		Path:   "/ecommerce/v1/orders/:id",
	})
}

func assertRoutePresent(t *testing.T, gotRoutes gin.RoutesInfo, wantRoute gin.RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			return
		}
	}
	t.Errorf("route %s %s not found", wantRoute.Method, wantRoute.Path)
}

func assertRouteNotPresent(t *testing.T, gotRoutes gin.RoutesInfo, wantRoute gin.RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			t.Errorf("route %s %s found", wantRoute.Method, wantRoute.Path)
		}
	}
}
