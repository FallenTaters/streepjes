package orderdomain

import "fmt"

type Price int

func (p Price) String() string {
	return fmt.Sprintf(`€%.2f`, float64(p)/100)
}

func (p Price) Times(n int) Price {
	return p * Price(n)
}
