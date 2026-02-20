package orderdomain_test

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestItemPrice(t *testing.T) {
	t.Parallel()

	item := orderdomain.Item{
		PriceGladiators: 200,
		PriceParabool:   150,
		PriceCalamari:   175,
	}

	t.Run("gladiators", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(200), item.Price(domain.ClubGladiators))
	})

	t.Run("parabool", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(150), item.Price(domain.ClubParabool))
	})

	t.Run("calamari", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(175), item.Price(domain.ClubCalamari))
	})

	t.Run("unknown club returns zero", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq(orderdomain.Price(0), item.Price(domain.ClubUnknown))
	})
}

func TestLinePrice(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	line := orderdomain.Line{
		Item:   orderdomain.Item{PriceGladiators: 200},
		Amount: 3,
	}
	assert.Eq(orderdomain.Price(600), line.Price(domain.ClubGladiators))
}
