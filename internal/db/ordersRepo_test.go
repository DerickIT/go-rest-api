package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/models/data"
	"github.com/derickit/go-rest-api/internal/util"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewOrderDataService(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	assert.Implements(t, (*db.OrdersDataService)(nil), dSvc)
}

func TestOrdersRepo_Create_Success(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}

	po := &data.Order{
		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Name(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}
	resultID, err := dSvc.Create(context.TODO(), po)
	if err != nil {
		t.Fail()
	}
	assert.Greater(t, len(resultID), 5)
}

func TestOrderRepo_Create_BadOrderID(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{
		ID: primitive.ObjectID{
			0x01, 0x02, 0x03, 0x04,
		},
		Version:  1,
		Products: products,
	}

	_, err := dSvc.Create(context.TODO(), po)
	require.Error(t, err)
	assert.EqualError(t, err, db.ErrInvalidPOIDCreate.Error())
}
func TestOrdersRepo_GetByIDSuccess(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{
		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}

	resultID, _ := dSvc.Create(context.TODO(), po)
	orderID, _ := primitive.ObjectIDFromHex(resultID)
	result, _ := dSvc.GetByID(context.TODO(), orderID)
	assert.NotNil(t, result)
	assert.EqualValues(t, orderID, result.ID)
}

func TestOrdersRepo_GetByIDSuccess_NoData(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	orderID, _ := primitive.ObjectIDFromHex("non-existent-id")
	result, err := dSvc.GetByID(context.TODO(), orderID)
	assert.Nil(t, result)
	assert.EqualError(t, err, db.ErrPOIDNotFound.Error())
}

func TestOrdersRepo_DeleteByIDSuccess(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{

		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}

	resultID, _ := dSvc.Create(context.TODO(), po)
	orderID, _ := primitive.ObjectIDFromHex(resultID)
	err := dSvc.DeleteByID(context.TODO(), orderID)
	require.NoError(t, err)
}

func TestOrdersRepo_DeleteByIDSuccess_NoData(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	orderID, _ := primitive.ObjectIDFromHex("non-existent-id")
	err := dSvc.DeleteByID(context.TODO(), orderID)
	assert.EqualError(t, err, db.ErrPOIDNotFound.Error())
}

func TestOrdersRepo_GetAll(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	results, _ := dSvc.GetAll(context.TODO(), int64(4))
	assert.Len(t, *results, 4)
}
func TestOrdersRepo_UpdateOrdersSucess(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{
		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}
	resultID, _ := dSvc.Create(context.TODO(), po)
	orderID, _ := primitive.ObjectIDFromHex(resultID)
	po.ID = orderID
	po.Status = data.OrderDelivered
	err := dSvc.Update(context.TODO(), po)
	require.NoError(t, err)
}

func TestOrdersRepo_UpdateOrder_EmptyOrderID(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{
		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}
	po.Status = data.OrderDelivered
	err := dSvc.Update(context.TODO(), po)
	assert.EqualError(t, err, db.ErrInvalidPOIDUpdate.Error())
}

func TestOrdersRepo_UpdateOrder_BadOrderID(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	po := &data.Order{
		ID:          primitive.ObjectID{},
		Version:     1,
		Products:    products,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}
	po.Status = data.OrderDelivered
	err := dSvc.Update(context.TODO(), po)
	assert.EqualError(t, err, db.ErrInvalidPOIDUpdate.Error())
}

func TestOrdersRepo_UpdateOrder_NonExistingID(t *testing.T) {
	d := testDBMgr.Database()
	dSvc := db.NewOrderRepo(d, lgr)
	products := []data.Product{
		{
			Name:     faker.Name(),
			Price:    util.RandomPrice(),
			UpdateAt: time.Now(),
		},
	}
	orderID, _ := primitive.ObjectIDFromHex("non-existent-id")
	po := &data.Order{
		ID:          orderID,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		User:        faker.Email(),
		Status:      data.OrderPending,
		TotalAmount: util.CalculateTotalAmount(products),
	}

	po.Status = data.OrderDelivered
	err := dSvc.Update(context.TODO(), po)
	assert.EqualError(t, err, db.ErrInvalidPOIDUpdate.Error())
}
