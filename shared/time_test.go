package shared

import (
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
)

func TestSameDate(t *testing.T) {
	t.Parallel()

	t.Run("same date same time", func(t *testing.T) {
		assert := assert.New(t)
		a := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
		b := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
		assert.True(SameDate(a, b))
	})

	t.Run("same date different time", func(t *testing.T) {
		assert := assert.New(t)
		a := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
		b := time.Date(2025, 6, 15, 23, 59, 59, 0, time.UTC)
		assert.True(SameDate(a, b))
	})

	t.Run("different dates", func(t *testing.T) {
		assert := assert.New(t)
		a := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
		b := time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC)
		assert.False(SameDate(a, b))
	})

	t.Run("different months", func(t *testing.T) {
		assert := assert.New(t)
		a := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
		b := time.Date(2025, 7, 15, 10, 0, 0, 0, time.UTC)
		assert.False(SameDate(a, b))
	})

	t.Run("different years", func(t *testing.T) {
		assert := assert.New(t)
		a := time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC)
		b := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
		assert.False(SameDate(a, b))
	})
}

func TestPrettyDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC)

	t.Run("today", func(t *testing.T) {
		assert := assert.New(t)
		assert.Eq("Today", prettyDate(now, now))
	})

	t.Run("yesterday", func(t *testing.T) {
		assert := assert.New(t)
		yesterday := now.AddDate(0, 0, -1)
		assert.Eq("Yesterday", prettyDate(yesterday, now))
	})

	t.Run("tomorrow", func(t *testing.T) {
		assert := assert.New(t)
		tomorrow := now.AddDate(0, 0, 1)
		assert.Eq("Tomorrow", prettyDate(tomorrow, now))
	})

	t.Run("same year", func(t *testing.T) {
		assert := assert.New(t)
		date := time.Date(2025, 3, 5, 10, 0, 0, 0, time.UTC)
		assert.Eq("5 Mar", prettyDate(date, now))
	})

	t.Run("different year", func(t *testing.T) {
		assert := assert.New(t)
		date := time.Date(2024, 3, 5, 10, 0, 0, 0, time.UTC)
		assert.Eq("5 Mar 2024", prettyDate(date, now))
	})
}

func TestPrettyTime(t *testing.T) {
	t.Parallel()

	t.Run("just now", func(t *testing.T) {
		assert := assert.New(t)
		now := time.Date(2025, 6, 15, 14, 0, 30, 0, time.UTC)
		target := time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC)
		assert.Eq("Just now", prettyTime(target, now))
	})

	t.Run("over a minute ago", func(t *testing.T) {
		assert := assert.New(t)
		now := time.Date(2025, 6, 15, 14, 5, 0, 0, time.UTC)
		target := time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC)
		assert.Eq("14:00", prettyTime(target, now))
	})

	t.Run("future time is not just now", func(t *testing.T) {
		assert := assert.New(t)
		now := time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC)
		target := time.Date(2025, 6, 15, 14, 0, 30, 0, time.UTC)
		assert.Eq("14:00", prettyTime(target, now))
	})
}

func TestPrettyDatetime(t *testing.T) {
	t.Parallel()

	t.Run("just now returns only just now", func(t *testing.T) {
		assert := assert.New(t)
		result := PrettyDatetime(time.Now().Add(-10 * time.Second))
		assert.Eq("Just now", result)
	})

	t.Run("combines date and time", func(t *testing.T) {
		assert := assert.New(t)
		target := time.Now().Add(-5 * time.Minute)
		result := PrettyDatetime(target)
		assert.Eq("Today "+target.Format("15:04"), result)
	})
}
