package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	errors2 "github.com/derickit/go-rest-api/internal/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/derickit/go-rest-api/internal/db/mocks"
	"github.com/derickit/go-rest-api/internal/handlers"
	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models"
	"github.com/derickit/go-rest-api/internal/models/data"
	"github.com/derickit/go-rest-api/internal/models/external"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func UnMarshalOrderData(d []byte) (*data.Order, error) {
	var r data.Order
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func UnMarshalOrdersData(d []byte) (*[]data.Order, error) {
	var r []data.Order
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func TestOrdersHandler_Create_Success(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "1", nil
		},
	}, lgr)

	r.POST("/orders", handler.Create)

	orderInput := external.OrderInput{

		Products: []external.ProductInput{
			{Name: "product 1",
				Price:    10.0,
				Quantity: 1},
		},
	}
	body, _ := json.Marshal(orderInput)
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var responseOrder external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &responseOrder)
	require.NoError(t, err)
	assert.Equal(t, int64(1), responseOrder.Version)
	assert.NotNil(t, responseOrder.CreatedAt)
	assert.NotNil(t, responseOrder.UpdatedAt)
	assert.Equal(t, orderInput.Products[0].Name, responseOrder.Products[0].Name)
	assert.InEpsilon(t, orderInput.Products[0].Name, responseOrder.Products[0].Price, 0)
	assert.Equal(t, orderInput.Products[0].Quantity, responseOrder.Products[0].Quantity)
	assert.InEpsilon(t, 20.0, responseOrder.TotalAmount, 0)
	assert.Equal(t, data.OrderPending, responseOrder.Status)
}

func TestOrderHandler_Create_InvalidInput(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "MOCK_ORDER_ID", nil
		},
	}, lgr)
	r.POST("/orders", handler.Create)
	invalidInput := "{invalid JSON}"
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(invalidInput)))

	r.ServeHTTP(recorder, c.Request)
	var apiErr external.APIError
	err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatusCode)
	assert.Equal(t, "orders_create_invalid_input", apiErr.ErrorCode)
	assert.Equal(t, "invalid order request body", apiErr.Message)
}

func TestOrdersHandler_Create_InternalServerError(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "", assert.AnError
		},
	}, lgr)
	r.POST("/orders", handler.Create)
	orderInput := external.OrderInput{
		Products: []external.ProductInput{
			{Name: "product 1",
				Price:    10.0,
				Quantity: 3},
		},
	}
	body, _ := json.Marshal(orderInput)
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var apiErr external.APIError
	err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiErr.HTTPStatusCode)
	assert.Equal(t, errors2.UnexpectedErrorMessage, apiErr.Message)
}

func TestGetAllOrdersSuccess(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetAllFunc: func(ctx context.Context, limit int64) (*[]data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/orders.json")
			if err != nil {
				return nil, err
			}
			dataOrders, _ := UnMarshalOrdersData(dataBytes)
			return dataOrders, nil
		},
	}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respOrders []external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &respOrders)
	require.NoError(t, err)
	assert.Len(t, respOrders, 10)
}

func TestGetAllOrdersFailure_DBRead(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetAllFunc: func(ctx context.Context, limit int64) (*[]data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/non-existent.json")
			if err != nil {
				return nil, err
			}
			dataOrders, _ := UnMarshalOrdersData(dataBytes)
			return dataOrders, nil
		},
	}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetAllOrdersFailure_LimitOutOfBounds(t *testing.T) {
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "10000")
	c.Request.URL.RawQuery = q.Encode()
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetAllOrdersFailure_InvalidLimit(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "aaabbbccc")
	c.Request.URL.RawQuery = q.Encode()
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

}

func TestGetOrderByIDSuccess(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/order.json")
			if err != nil {
				return nil, err
			}
			dataOrder, _ := UnMarshalOrderData(dataBytes)
			return dataOrder, nil
		},
	}, lgr)
	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/1", nil)

	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	var respOrders external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &respOrders)
	require.NoError(t, err)
	assert.Equal(t, "1", respOrders.ID)

}

func TestGetOrderByID_DBReadFailure(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
			return nil, errors.New("db error")
		},
	}, lgr)

	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/1", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetOrderByID_BadPathParam(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
			return nil, errors.New("db error")
		},
	}, lgr)
	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/''", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteOrderByIDSuccess(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(_ context.Context, _ primitive.ObjectID) error {
			return nil
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/orders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/1", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteOrderByID_DBFailure(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(_ context.Context, _ primitive.ObjectID) error {
			return errors.New("db error")
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/roders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/1", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeleteOrderByID_BadPathParam(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(ctx context.Context, id primitive.ObjectID) error {
			return nil
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/orders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/''", nil)
	r.ServeHTTP(recorder, c.Request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
