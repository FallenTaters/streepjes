package auth

import (
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
)

type Activewatcher struct {
	Show bool `vugu:"data"`
}

func (aw *Activewatcher) Init() {
	events.Listen(events.InactiveWarning, `activewatcher`, func() {
		defer global.LockAndRender()()
		aw.Show = true
	})

	events.Listen(events.Unauthorized, `activewatcher`, func() {
		defer global.LockAndRender()()
		aw.Show = false
	})
}

func (aw *Activewatcher) click() {
	aw.Show = false
}
