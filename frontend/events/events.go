package events

type Event int

const (
	Unauthorized Event = iota + 1
	Login
)

var listeners = make(map[Event]map[string]func())

func Trigger(event Event) {
	for _, listener := range listeners[event] {
		if listener != nil {
			go listener()
		}
	}
}

func Listen(event Event, key string, callback func()) {
	if listeners[event] == nil {
		listeners[event] = make(map[string]func())
	}

	listeners[event][key] = callback
}
