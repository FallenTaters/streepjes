package beercss

import (
	"github.com/PotatoesFall/vecty-test/frontend/jscall/beercss"
	"github.com/vugu/vugu"
)

type Input struct {
	Label string `vugu:"data"`
}

func (i *Input) Rendered(vugu.RenderedCtx) {
	go beercss.UI()
}
