package orderdomain

import "github.com/PotatoesFall/vecty-test/domain"

type Member struct {
	ID   int         `json:"id"`
	Club domain.Club `json:"club"`
	Name string      `json:"name"`
}
