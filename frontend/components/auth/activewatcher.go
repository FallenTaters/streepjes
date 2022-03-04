package auth

import (
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
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
