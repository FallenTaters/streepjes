package shared

type Club int

const (
	ClubUnknown    Club = iota // Unknown
	ClubParabool               // Parabool
	ClubGladiators             // Gladiators
)

func (c Club) Int() int {
	return int(c)
}

//go:generate enumer -linecomment -type Club
