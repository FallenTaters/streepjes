package orderdomain

import (
	"time"

	"github.com/PotatoesFall/vecty-test/domain"
)

//go:generate enumer -json -sql -linecomment -type Status

type Order struct {
	ID          int         `json:"id"`
	Club        domain.Club `json:"club"`
	BartenderID int         `json:"bartender"`
	MemberID    int         `json:"memberId"`
	Contents    string      `json:"contents"`
	Price       Price       `json:"price"`
	OrderTime   time.Time   `json:"orderDate"`
	Status      Status      `json:"status"`
	StatusTime  time.Time   `json:"statusDate"`
}

type Status int

const (
	StatusOpen      Status = iota + 1 // Open
	StatusBilled                      // Billed
	StatusPaid                        // Paid
	StatusCancelled                   // Cancelled
)
