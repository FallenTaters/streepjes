package orderdomain

import (
	"time"

	"github.com/FallenTaters/streepjes/domain"
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

type Month struct {
	Year  int
	Month time.Month
}

func MonthOf(t time.Time) Month {
	return Month{
		t.Year(),
		t.Month(),
	}
}

func (m Month) Time() time.Time {
	return time.Date(m.Year, m.Month, 1, 0, 0, 0, 0, time.Local)
}
