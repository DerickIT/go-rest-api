package util

import (
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/derickit/go-rest-api/internal/models/data"
)

func FormatTimeToISO(timeToFormat time.Time) string {
	return timeToFormat.Format(time.RFC3339)
}

func CurrentISOTime() string {
	return FormatTimeToISO(time.Now().UTC())
}

func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}

const defaultPrice = 100
const MaxPrice = 1000

func RandomPrice() float64 {
	var price *big.Int
	var err error
	if price, err = rand.Int(rand.Reader, big.NewInt(MaxPrice)); err != nil {
		price = big.NewInt(defaultPrice)
	}
	pf, _ := price.Float64()
	return pf
}

func CalculateTotalAmount(products []data.Product) float64 {
	var total float64
	for _, product := range products {
		total += product.Price * (float64(product.Quantity))
	}
	return total
}
