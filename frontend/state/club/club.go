package club

type ClubType int

const (
	Gladiators ClubType = iota + 1
	Parabool
)

var Club ClubType

var Listeners []func(ClubType)

func Toggle() {
	if Club == Gladiators {
		Club = Parabool
	} else {
		Club = Gladiators
	}

	for _, f := range Listeners {
		f(Club)
	}
}
