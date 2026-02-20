package orderdomain

import (
	"fmt"
	"time"

	"github.com/FallenTaters/streepjes/domain"
)

//go:generate go tool enumer -json -sql -linecomment -type Status

type Order struct {
	ID          int         `json:"id"`
	Club        domain.Club `json:"club"`
	BartenderID int         `json:"bartender"`
	MemberID    int         `json:"memberId"`
	Contents    string      `json:"contents"` // Usually []Orderline as JSON
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

func ParseMonth(s string) (Month, error) {
	t, err := time.Parse(`2006-01`, s)
	return MonthOf(t), err
}

func CurrentMonth() Month {
	return MonthOf(time.Now())
}

func (m Month) Start() time.Time {
	return time.Date(m.Year, m.Month, 1, 0, 0, 0, 0, time.UTC)
}

func (m Month) End() time.Time {
	return time.Date(m.Year, m.Month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
}

func (m Month) String() string {
	return fmt.Sprintf(`%04d-%02d`, m.Year, m.Month)
}
