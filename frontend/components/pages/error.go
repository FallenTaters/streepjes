package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Error(msg string) vecty.Component {
	return &err{msg: msg}
}

type err struct {
	vecty.Core

	msg string
}

func (err *err) Render() vecty.ComponentOrHTML {
	return elem.Article(
		vecty.Markup(vecty.Class(`error`)),
		vecty.Text(err.msg),
	)
}
