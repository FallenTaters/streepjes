package order

// import (
// 	"github.com/PotatoesFall/vecty-test/frontend/store"
// 	"github.com/hexops/vecty"
// 	"github.com/hexops/vecty/elem"
// )

// type Summary struct {
// 	vecty.Core
// }

// func (o *Summary) Render() vecty.ComponentOrHTML {
// 	return elem.Div(
// 		vecty.Markup(vecty.Class(`row`, `no-wrap`)),
// 		elem.Heading3(
// 			vecty.Markup(vecty.Class(`col`, `min`)),
// 			vecty.Text(`Total`),
// 		),
// 		elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
// 		elem.Heading3(
// 			vecty.Markup(vecty.Class(`col`, `min`)),
// 			vecty.Text(store.Order.CalculateTotal().String()),
// 		),
// 	)
// }
