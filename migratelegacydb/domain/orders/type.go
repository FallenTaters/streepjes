package orders

import (
	"time"

	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/null"
)

type OrderStatus int

const (
	OrderStatusOpen      OrderStatus = iota + 1 // Open
	OrderStatusBilled                           // Billed
	OrderStatusPaid                             // Paid
	OrderStatusCancelled                        // Cancelled
)

//go:generate enumer -json -linecomment -type OrderStatus

type Order struct {
	ID         int         `json:"id"`
	Club       shared.Club `json:"club"`
	Bartender  string      `json:"bartender"`
	MemberID   int         `json:"memberId"`
	Contents   string      `json:"contents"`
	Price      int         `json:"price"`
	OrderTime  time.Time   `json:"orderDate"`
	Status     OrderStatus `json:"status"`
	StatusTime time.Time   `json:"statusDate"`
}

func (o Order) Key() []byte {
	return buckets.Itob(o.ID)
}

type OrderFilter struct {
	Club      null.Int
	Bartender null.String
	Member    null.Int
	Status    []int
	Month     *time.Time
}
