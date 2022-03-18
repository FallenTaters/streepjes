package orderdomain

import (
	"time"

	"github.com/FallenTaters/streepjes/domain"
)

type Member struct {
	ID        int         `json:"id"`
	Club      domain.Club `json:"club"`
	Name      string      `json:"name"`
	LastOrder time.Time   `json:"last_order"`
}
