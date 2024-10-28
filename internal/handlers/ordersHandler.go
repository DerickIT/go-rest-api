package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/errors"
	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models/data"
	"github.com/derickit/go-rest-api/internal/models/external"
	"github.com/derickit/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OrderIDPath = "id"
	MaxPageSize = 100
)

type OrdersHandler struct {
	oDataSvc db.OrdersDataService
	logger   *logger.AppLogger
}

func NewOrdersHandler(dSvc db.OrdersDataService, lgr *logger.AppLogger) *OrdersHandler {
	o := &OrdersHandler{
		oDataSvc: dSvc,
		logger:   lgr,
	}
	return o
}

func (o *OrdersHandler) Create(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	var orderInput external.OrderInput
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		apiErr := &external.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      errors.OrderCreateInvalidInput,
			Message:        "Invalid order request body",
			DebugID:        requestID,
		}
		lgr.Error().Err(err).Int("HttpStatusCode", apiErr.HTTPStatusCode).Str("ErrorCode", apiErr.ErrorCode).Msg(apiErr.Message)
		c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
		return
	}
	var products []data.Product
	for _, productInput := range orderInput.Products {
		product := data.Product{
			Name:     productInput.Name,
			Price:    productInput.Price,
			Quantity: productInput.Quantity,
		}
		products = append(products, product)
	}

	order := data.Order{
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Products:    products,
		User:        faker.Email(),
		TotalAmount: util.CalculateTotalAmount(products),
		Status:      data.OrderPending,
	}

	if id, err := o.oDataSvc.Create(c, &order); err == nil {
		extOrder := external.Order{
			ID:          id,
			Version:     order.Version,
			CreatedAt:   util.FormatTimeToISO(order.CreatedAt),
			UpdatedAt:   util.FormatTimeToISO(order.UpdatedAt),
			Products:    order.Products,
			User:        order.User,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
		}
		c.JSON(http.StatusCreated, extOrder)
		return
	}

	apiErr := &external.APIError{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorCode:      errors.OrderCreateServerError,
		Message:        errors.UnexpectedErrorMessage,
		DebugID:        requestID,
	}

	lgr.Error().Int("HttpStatusCode", apiErr.HTTPStatusCode).Str("ErrorCode", apiErr.ErrorCode).Msg(apiErr.Message)
	c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)

}

func (o *OrdersHandler) GetAll(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	limit, apiErr := o.parseLimitQueryParam(c)
	// Validate  inputs : fail fast order
	// Parse query params
	if apiErr != nil {
		c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
		return
	}

	orders, err := o.oDataSvc.GetAll(c, limit)
	var extOrders []external.Order
	if orders != nil {
		extOrders = make([]external.Order, len(*orders))
		for i, o := range *orders {
			extOrders[i] = external.Order{
				ID:          o.ID.Hex(),
				Version:     o.Version,
				Status:      o.Status,
				TotalAmount: o.TotalAmount,
				User:        o.User,
				CreatedAt:   util.FormatTimeToISO(o.CreatedAt),
				UpdatedAt:   util.FormatTimeToISO(o.UpdatedAt),
				Products:    o.Products,
			}
		}
	}

	if err != nil {
		aErr := &external.APIError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      errors.OrderGetServerError,
			Message:        errors.UnexpectedErrorMessage,
			DebugID:        requestID,
		}
		lgr.Error().
			Int("HttpStatusCode", aErr.HTTPStatusCode).
			Str("ErrorCode", aErr.ErrorCode).
			Msg(aErr.Message)
		c.AbortWithStatusJSON(aErr.HTTPStatusCode, aErr)
		return
	}
	c.JSON(http.StatusOK, extOrders)

}

func (o *OrdersHandler) GetByID(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	id := c.Param(OrderIDPath)
	oID, err := primitive.ObjectIDFromHex(id)
	if oID.IsZero() || err != nil {
		aErr := &external.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "",
			Message:        "Invalid order id",
			DebugID:        requestID,
		}
		lgr.Error().Int("HttpStatusCode", aErr.HTTPStatusCode).Str("ErrorCode", aErr.ErrorCode).Msg(aErr.Message)
		return
	}
	order, err := o.oDataSvc.GetByID(c, oID)
	if err != nil {
		aErr := &external.APIError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "",
			Message:        "fill this in with a meaningful error message",
			DebugID:        requestID,
		}
		lgr.Error().Int("HttpStatusCode", aErr.HTTPStatusCode).Str("ErrorCode", aErr.ErrorCode).Msg(aErr.Message)
		c.AbortWithStatusJSON(aErr.HTTPStatusCode, aErr)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrdersHandler) parseLimitQueryParam(c *gin.Context) (int64, *external.APIError) {
	lgr, requestID := o.logger.WithReqID(c)
	l := db.DefaultPageSize
	if input, exists := c.GetQuery("limit"); exists && input != "" {
		var err error
		l, err = strconv.Atoi(input)
		if err != nil {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        fmt.Sprintf("Integer value within 1 and %d is expected for limit query param", MaxPageSize),
				DebugID:        requestID,
			}
			lgr.Error().Int("HttpStatusCode", apiErr.HTTPStatusCode).Str("ErrorCode", apiErr.ErrorCode).Msg(apiErr.Message)
			return 0, apiErr
		}
		if l < 1 || l > MaxPageSize {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        fmt.Sprintf("Integer value within 1 and %d is expected for limit query param", MaxPageSize),
				DebugID:        requestID,
			}
			lgr.Error().Int("HttpStatusCode", apiErr.HTTPStatusCode).Str("ErrorCode", apiErr.ErrorCode).Msg(apiErr.Message)
			return 0, apiErr
		}

	}
	return int64(1), nil
}
