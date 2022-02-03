package domain

//go:generate enumer -json -linecomment -type Club

type Club int

const (
	ClubUnknown    Club = iota // Unknown
	ClubParabool               // Parabool
	ClubGladiators             // Gladiators
)
