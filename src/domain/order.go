package domain

import "time"

//go:generate enumer -json -linecomment -type OrderStatus

type Order struct {
	ID         int         `json:"id"`
	Club       Club        `json:"club"`
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
