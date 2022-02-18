package global

import "github.com/vugu/vugu"

var EventEnv vugu.EventEnv

func LockAndRender() func() {
	EventEnv.Lock()
	return func() {
		EventEnv.UnlockRender()
	}
}

func LockOnly() func() {
	EventEnv.Lock()
	return func() {
		EventEnv.UnlockOnly()
	}
}
