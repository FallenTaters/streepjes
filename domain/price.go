package domain

import "fmt"

type Price int

func (p Price) String() string {
	return fmt.Sprintf(`€%.2f`, float64(p)/100)
}
