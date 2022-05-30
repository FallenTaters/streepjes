package orderdomain

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Price int

var printer = message.NewPrinter(language.English)

func (p Price) String() string {
	return printer.Sprintf("€%.2f\n", float64(p)/100)
}

func (p Price) Times(n int) Price {
	return p * Price(n)
}
