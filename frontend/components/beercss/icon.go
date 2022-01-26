package beercss

import "github.com/hexops/vecty"

type IconType string

const (
	IconTypeLocalBar  IconType = `local_bar`
	IconTypeHistory   IconType = `history`
	IconTypeAddCircle IconType = `add_circle`
	IconTypeList      IconType = `list`
	IconTypePerson    IconType = `person`
	IconTypeSwapHoriz IconType = `swap_horiz`
	IconTypePayments  IconType = `payments`
)

func Icon(i IconType, markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Text(string(i)))
	return vecty.Tag(`i`, markup...)
}
