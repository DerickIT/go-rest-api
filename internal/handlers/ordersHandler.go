package handlers

import (
	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/logger"
)

const (
	OrderIDPath = "id"
	MaxPageSize = 100
)

type OrdersHandler struct {
	oDataSvc db.OrdersDataService
	logger   *logger.AppLogger
}
