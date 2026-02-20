package orderdomain_test

import (
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestParseMonth(t *testing.T) {
	t.Parallel()

	t.Run("valid month", func(t *testing.T) {
		assert := assert.New(t)
		m, err := orderdomain.ParseMonth("2025-06")
		assert.NoError(err)
		assert.Eq(2025, m.Year)
		assert.Eq(time.June, m.Month)
	})

	t.Run("january", func(t *testing.T) {
		assert := assert.New(t)
		m, err := orderdomain.ParseMonth("2025-01")
		assert.NoError(err)
		assert.Eq(time.January, m.Month)
	})

	t.Run("december", func(t *testing.T) {
		assert := assert.New(t)
		m, err := orderdomain.ParseMonth("2025-12")
		assert.NoError(err)
		assert.Eq(time.December, m.Month)
	})

	t.Run("invalid format", func(t *testing.T) {
		assert := assert.New(t)
		_, err := orderdomain.ParseMonth("not-a-month")
		assert.Error(err)
	})

	t.Run("empty string", func(t *testing.T) {
		assert := assert.New(t)
		_, err := orderdomain.ParseMonth("")
		assert.Error(err)
	})
}

func TestMonthOf(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	m := orderdomain.MonthOf(time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC))
	assert.Eq(2025, m.Year)
	assert.Eq(time.March, m.Month)
}

func TestMonthStart(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	m := orderdomain.Month{Year: 2025, Month: time.June}
	start := m.Start()
	assert.Eq(time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC), start)
}

func TestMonthEnd(t *testing.T) {
	t.Parallel()

	t.Run("regular month", func(t *testing.T) {
		assert := assert.New(t)
		m := orderdomain.Month{Year: 2025, Month: time.June}
		assert.Eq(time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC), m.End())
	})

	t.Run("december rolls to next year", func(t *testing.T) {
		assert := assert.New(t)
		m := orderdomain.Month{Year: 2025, Month: time.December}
		assert.Eq(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), m.End())
	})
}

func TestMonthString(t *testing.T) {
	t.Parallel()

	t.Run("regular month", func(t *testing.T) {
		assert := assert.New(t)
		m := orderdomain.Month{Year: 2025, Month: time.June}
		assert.Eq("2025-06", m.String())
	})

	t.Run("january with padding", func(t *testing.T) {
		assert := assert.New(t)
		m := orderdomain.Month{Year: 2025, Month: time.January}
		assert.Eq("2025-01", m.String())
	})

	t.Run("december", func(t *testing.T) {
		assert := assert.New(t)
		m := orderdomain.Month{Year: 2025, Month: time.December}
		assert.Eq("2025-12", m.String())
	})
}
