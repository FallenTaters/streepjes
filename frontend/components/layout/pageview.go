package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/vugu/vugu"
)

type Page int

const (
	PageLogin Page = iota + 1

	PageOrder
	PageHistory

	PageCatalog
	PageMembers
	PageUsers
)

type Pageview struct {
	Page Page
}

func (pv *Pageview) Init(vugu.InitCtx) {
	events.Listen(events.Unauthorized, `pageview`, func() {
		global.EventEnv.Lock()
		defer global.EventEnv.UnlockRender()

		pv.Page = PageLogin
	})

	pv.Page = PageOrder
}

func getStyles() string {
	switch window.GetSize() {
	case window.SizeL:
		return `padding: 20px 40px 20px 120px;`

	case window.SizeM:
		return `padding: 20px 20px 0px 100px;`

	case window.SizeS:
		return `padding: 20px 10px 50px 10px;`
	}

	return ``
}
