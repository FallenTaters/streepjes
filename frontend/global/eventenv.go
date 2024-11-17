package global

import (
	"time"

	"github.com/vugu/vugu"
)

var EventEnv vugu.EventEnv

func LockAndRender() func() {
	time.Sleep(10 * time.Millisecond)

	EventEnv.Lock()
	return func() {
		EventEnv.UnlockRender()
	}
}

func LockOnly() func() {
	time.Sleep(10 * time.Millisecond)

	EventEnv.Lock()
	return func() {
		EventEnv.UnlockOnly()
	}
}
