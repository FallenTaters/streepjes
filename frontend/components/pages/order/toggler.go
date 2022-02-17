package order

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/store"
	"github.com/vugu/vugu"
)

type Toggler struct {
	Club domain.Club `vugu:"data"`
}

func (t *Toggler) Compute(vugu.ComputeCtx) {
	t.Club = store.Order.Club
}

func (t *Toggler) toggle() {
	store.Order.ToggleClub()
}
