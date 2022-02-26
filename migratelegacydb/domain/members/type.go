package members

import (
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"
)

type Member struct {
	ID   int         `json:"id"`
	Club shared.Club `json:"club"`
	Name string      `json:"name"`
	Debt int         `json:"debt"`
}

func (m Member) Key() []byte {
	return buckets.Itob(m.ID)
}
