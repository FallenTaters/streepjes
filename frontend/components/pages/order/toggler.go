package order

import (
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/frontend/store"
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
