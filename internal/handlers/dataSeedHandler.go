package handlers

import (
	"net/http"
	"time"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/models/data"
	"github.com/derickit/go-rest-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

const (
	seedRecoedCount = 10000
)

type SeedHandler struct {
	oDataSvc db.OrdersDataService
}

func NewDataSeedHandler(svc db.OrdersDataService) *SeedHandler {
	sc := &SeedHandler{
		oDataSvc: svc,
	}
	return sc
}

func (s *SeedHandler) SeedDB(c *gin.Context) {
	for i := 0; i < seedRecoedCount; i++ {
		products := []data.Product{
			{
				Name:     faker.Name(),
				Price:    util.RandomPrice(),
				UpdateAt: time.Now(),
			},
			{
				Name:     faker.Name(),
				Price:    util.RandomPrice(),
				UpdateAt: time.Now(),
			},
		}
		po := &data.Order{
			Version:     1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Products:    products,
			User:        faker.Email(),
			Status:      data.OrderPending,
			TotalAmount: util.CalculateTotalAmount(products),
		}

		_, err := s.oDataSvc.Create(c, po)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Message": "Failed to seed data",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Data seeded successfully",
		"Count":   seedRecoedCount,
	})
}
