package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderSuccess(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 10.0,
		Tax:   2.0,
	}
	assert.Equal(t, order.ID, "123")
	assert.Equal(t, order.Price, 10.0)
	assert.Equal(t, order.Tax, 2.0)
	assert.Nil(t, order.IsValid())
}

func TestCreateOrderSuccessWithFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	assert.Equal(t, order.ID, "123")
	assert.Equal(t, order.Price, 10.0)
	assert.Equal(t, order.Tax, 2.0)

	err = order.CalculateFinalPrice()
	assert.Nil(t, err)
	assert.Equal(t, order.FinalPrice, 12.0)
}

func TestErrorWhenCreateOrderWithouID(t *testing.T) {
	order := Order{}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestErrorWhenCreateOrderWithoutPrice(t *testing.T) {
	order := Order{
		ID: "123",
	}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestErrorWhenCreateOrderWithouTax(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 2.0,
	}
	assert.Error(t, order.IsValid(), "invalid id")
}
