package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
)

type Header struct {
	Navigate func(Page)

	Size window.Size
}

func (h *Header) Init() {
	h.Size = window.OnResize(`header`, func(s window.Size) {
		h.Size = s
	})
}

func (h *Header) menuClasses() string {
	side := `left`
	if h.Size == window.SizeS {
		side = `bottom`
	}

	return `menu ` + side
}
