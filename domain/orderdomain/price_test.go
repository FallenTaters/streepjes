package orderdomain_test

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestPriceString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		price orderdomain.Price
		want  string
	}{
		{"zero", 0, "€0.00"},
		{"one euro", 100, "€1.00"},
		{"cents only", 50, "€0.50"},
		{"mixed", 350, "€3.50"},
		{"large amount", 123456, "€1,234.56"},
		{"single cent", 1, "€0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Eq(tt.want, tt.price.String())
		})
	}
}

func TestPriceTimes(t *testing.T) {
	t.Parallel()

	t.Run("multiply by positive", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(600), orderdomain.Price(200).Times(3))
	})

	t.Run("multiply by zero", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(0), orderdomain.Price(200).Times(0))
	})

	t.Run("multiply by one", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(200), orderdomain.Price(200).Times(1))
	})
}
