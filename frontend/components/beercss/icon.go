package beercss

import "github.com/hexops/vecty"

type IconType string

const (
	IconTypeLocalBar IconType = `local_bar`
	IconTypeHistory  IconType = `history`
)

func Icon(i IconType, markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Text(string(i)))
	return vecty.Tag(`i`, markup...)
}
