package events

import "sync"

type Event int

const (
	Unauthorized Event = iota + 1
	Login

	InactiveWarning

	OrderPlaced
)

var (
	listeners = make(map[Event]map[string]func())
	mu        sync.RWMutex
)

func Trigger(event Event) {
	mu.RLock()
	defer mu.RUnlock()

	for _, listener := range listeners[event] {
		if listener != nil {
			go listener()
		}
	}
}

func Listen(event Event, key string, callback func()) {
	mu.Lock()
	defer mu.Unlock()

	if listeners[event] == nil {
		listeners[event] = make(map[string]func())
	}

	listeners[event][key] = callback
}
