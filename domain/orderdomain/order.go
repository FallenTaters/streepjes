package orderdomain

import (
	"time"

	"github.com/PotatoesFall/vecty-test/domain"
)

//go:generate enumer -json -linecomment -type OrderStatus

type Order struct {
	ID         int         `json:"id"`
	Club       domain.Club `json:"club"`
	Bartender  string      `json:"bartender"`
	MemberID   int         `json:"memberId"`
	Contents   string      `json:"contents"`
	Price      int         `json:"price"`
	OrderTime  time.Time   `json:"orderDate"`
	Status     OrderStatus `json:"status"`
	StatusTime time.Time   `json:"statusDate"`
}

type OrderStatus int

const (
	OrderStatusOpen      OrderStatus = iota + 1 // Open
	OrderStatusBilled                           // Billed
	OrderStatusPaid                             // Paid
	OrderStatusCancelled                        // Cancelled
)
