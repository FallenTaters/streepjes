package shared

import "time"

func PrettyDatetime(t time.Time) string {
	t, now := t.In(time.Local), time.Now()
	timePretty := prettyTime(t, now)
	if timePretty == `Just now` {
		return timePretty
	}

	return prettyDate(t, now) + ` ` + timePretty
}

func PrettyDate(t time.Time) string {
	return prettyDate(t.In(time.Local), time.Now())
}

func PrettyTime(t time.Time) string {
	return prettyTime(t.In(time.Local), time.Now())
}

func prettyDate(t, now time.Time) string {
	if SameDate(t, now) {
		return `Today`
	}
	if SameDate(t, now.AddDate(0, 0, 1)) {
		return `Tomorrow`
	}
	if SameDate(t, now.AddDate(0, 0, -1)) {
		return `Yesterday`
	}

	if now.Year() == t.Year() {
		t.Format(`2 Jan`)
	}

	return t.Format(`2 Jan 2006`)
}

func prettyTime(t, now time.Time) string {
	if now.After(t) && now.Sub(t) < time.Minute {
		return `Just now`
	}

	return t.Format(`03:04`)
}

func SameDate(a, b time.Time) bool {
	year1, month1, day1 := a.Year(), a.Month(), a.Day()
	year2, month2, day2 := b.Year(), b.Month(), b.Day()

	return year1 == year2 && month1 == month2 && day1 == day2
}
